// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package apievent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/5GSEC/SentryFlow/protobuf/golang"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	statusCode string
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter events by response status code (specific or family)",
	Long: `Filter events streamed from SentryFlow by response status code.

You can filter by specific status codes (e.g., 200, 404) to match only those codes,
or by status code families (e.g., 2xx, 4xx) to match a range of codes
(e.g., 2xx matches all codes between 200 and 299).`,
	Example: `
# Print filtered API Events based on some Response Status code
sfctl event filter --status "200"

# Print filtered API Events based on family of Response Status code
sfctl event filter --status "2xx"
sfctl event filter --status "3xx"
sfctl event filter --status "4xx"
sfctl event filter --status "5xx"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}
		if val == "" || len(val) < 3 ||
			!strings.HasPrefix(val, "1") &&
				!strings.HasPrefix(val, "2") &&
				!strings.HasPrefix(val, "3") &&
				!strings.HasPrefix(val, "4") &&
				!strings.HasPrefix(val, "5") {
			return fmt.Errorf("invalid status code, %v", val)
		}
		return filterEvents(cmd.Flags())
	},
	SilenceUsage: true,
}

func filterEvents(flags *pflag.FlagSet) error {
	statusCode, _ = flags.GetString("status")
	return printEvents(flags)
}

func printFilteredEvents(event *pb.APIEvent, prettyPrint bool) {
	var err error

	apiResStatusCodeStr, exists := event.Response.Headers[":status"]
	if !exists {
		return
	}

	apiResStatusCode, err := strconv.Atoi(apiResStatusCodeStr)
	if err != nil {
		logger.Warnf("failed to parse status code: %v", err)
		return
	}

	if !matches(apiResStatusCode, statusCode) {
		return
	}

	var jsonBytes []byte
	if prettyPrint {
		jsonBytes, err = json.MarshalIndent(event, "", "  ")
		if err != nil {
			logger.Warnf("failed to marshal event: %v", err)
			return
		}
	} else {
		jsonBytes, err = json.Marshal(event)
		if err != nil {
			logger.Warnf("failed to marshal event: %v", err)
		}
	}

	fmt.Println(string(jsonBytes))
}

func matches(curr int, target string) bool {
	switch target {
	case "2xx":
		return curr >= 200 && curr < 300
	case "3xx":
		return curr >= 300 && curr < 400
	case "4xx":
		return curr >= 400 && curr < 500
	case "5xx":
		return curr >= 500 && curr < 600
	default:
		if specificCode, err := strconv.Atoi(target); err == nil {
			return curr == specificCode
		}
		return false
	}
}

func init() {
	filterCmd.Flags().StringVar(&statusCode, "status", "", "response status code")
}
