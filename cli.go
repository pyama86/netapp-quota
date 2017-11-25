package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/johntdyer/slackrus"
	"github.com/pepabo/go-netapp/netapp"
	"github.com/sirupsen/logrus"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

}

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

const (
	statusQuotaOn  = "on"
	statusQuotaOff = "off"
)

var (
	version   string
	revision  string
	goversion string
	builddate string
	builduser string
)

func printVersion() {
	fmt.Printf("netapp-quota version: %s (%s)\n", version, revision)
	fmt.Printf("build at %s (with %s) by %s\n", builddate, goversion, builduser)
}

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

var Name = "netapp-quota"

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		url          string
		user         string
		password     string
		prefix       string
		svm          string
		onInterval   int
		offInterval  int
		slackURL     string
		slackChannel string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&url, "url", "", "netapp api endpoint")

	flags.StringVar(&user, "user", "", "netapp api BasicAuthUser")

	flags.StringVar(&password, "password", "", "netapp api BasicAuthPassword")

	flags.StringVar(&prefix, "prefix", "", "netapp volume prefix")

	flags.StringVar(&svm, "svm", "", "netapp svm server name")

	flags.IntVar(&onInterval, "on-interval", 10, "quota on interval")
	flags.IntVar(&offInterval, "off-interval", 300, "quota off interval")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	flags.StringVar(&slackURL, "slack-url", "", "slack webhook url")
	flags.StringVar(&slackChannel, "slack-channel", "", "slack channel")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		printVersion()
		return ExitCodeOK
	}

	if slackURL != "" && slackChannel != "" {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        slackURL,
			AcceptedLevels: slackrus.LevelThreshold(logrus.ErrorLevel),
			Channel:        slackChannel,
			IconEmoji:      ":ghost:",
			Username:       Name,
		})
	}

	client := netapp.NewClient(
		url,
		"1.20",
		&netapp.ClientOptions{
			BasicAuthUser:     user,
			BasicAuthPassword: password,
			SSLVerify:         false,
		},
	)

	go func() {
		for {
			fncSwitchQuota(quotaOff, client, svm, prefix, statusQuotaOn)
			time.Sleep(time.Duration(offInterval) * time.Second)
		}
	}()

	for {
		fncSwitchQuota(quotaOn, client, svm, prefix, statusQuotaOff)
		time.Sleep(time.Duration(onInterval) * time.Second)
	}
	return ExitCodeOK
}

func fncSwitchQuota(fn func(client *netapp.Client, vserver, volume string) error, client *netapp.Client, svm, prefix, workStatus string) {
	vls, err := volumeList(client, prefix)
	if err != nil {
		logrus.Error(err)
	}
	for _, vl := range vls {
		status, err := getQuotaStatus(client, svm, vl)
		if err != nil {
			logrus.Error(err)
		}

		if status == workStatus {
			ok, err := hasQuota(client, svm, vl)

			if err != nil {
				logrus.Error(err)
			}

			if !ok {
				continue
			}

			err = fn(client, svm, vl)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func hasQuota(client *netapp.Client, svm, volume string) (bool, error) {
	quotaList, _, err := client.Quota.List(&netapp.QuotaOptions{
		MaxRecords: 2,
		Query: &netapp.QuotaEntry{
			Vserver: svm,
			Volume:  volume,
		},
	})
	if err != nil {
		return false, err
	}

	if !quotaList.Results.Passed() {
		return false, fmt.Errorf("get quota list error")
	}

	for _, q := range quotaList.Results.AttributesList.QuotaEntry {
		if q.Volume == volume {
			return true, nil
		}
	}
	return false, nil

}
func volumeList(client *netapp.Client, prefix string) ([]string, error) {
	volumes := []string{}
	vl, _, err := client.Volume.List(nil)

	if err != nil {
		return nil, err
	}
	if !vl.Results.Passed() {
		return nil, fmt.Errorf("get volume list error")
	}

	for _, v := range vl.Results.AttributesList.VolumeAttributes {
		vn := v.VolumeIDAttributes.Name
		if prefix == "" || strings.HasPrefix(vn, prefix) {
			volumes = append(volumes, vn)
		}

	}
	return volumes, nil
}

func getQuotaStatus(client *netapp.Client, vserver, volume string) (string, error) {
	res, _, err := client.Quota.Status(vserver, volume)
	if err != nil {
		return "", err
	}

	if !res.Results.Passed() {
		return "", fmt.Errorf("getQuotaStatus failed: %s", res.Results.Reason)
	}

	return res.Results.QuotaStatus, nil
}

func switchQuota(res *netapp.QuotaStatusResponse, r *http.Response, err error) error {
	if err != nil {
		return err
	}
	if !res.Results.Passed() {
		return fmt.Errorf("switchQuota failed: %s", res.Results.Reason)
	}
	return nil
}

func quotaOn(client *netapp.Client, vserver, volume string) error {
	logrus.Infof("switch quota to on vserver=%s volume=%s", vserver, volume)
	return switchQuota(client.Quota.On(vserver, volume))
}

func quotaOff(client *netapp.Client, vserver, volume string) error {
	logrus.Infof("switch quota to off vserver=%s volume=%s", vserver, volume)
	return switchQuota(client.Quota.Off(vserver, volume))
}
