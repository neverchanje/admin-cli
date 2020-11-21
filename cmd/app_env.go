package cmd

import (
	"admin-cli/executor"
	"admin-cli/shell"
	"fmt"

	"github.com/desertbit/grumble"
)

// NOTE: some old-version servers may not support some of the keys.
var predefinedAppEnvKeys = []string{
	"rocksdb.usage_scenario",
	"replica.deny_client_write",
	"replica.write_throttling",
	"replica.write_throttling_by_size",
	"default_ttl",
	"manual_compact.disabled",
	"manual_compact.max_concurrent_running_count",
	"manual_compact.once.trigger_time",
	"manual_compact.once.target_level",
	"manual_compact.once.bottommost_level_compaction",
	"manual_compact.periodic.trigger_time",
	"manual_compact.periodic.target_level",
	"manual_compact.periodic.bottommost_level_compaction",
	"rocksdb.checkpoint.reserve_min_count",
	"rocksdb.checkpoint.reserve_time_seconds",
	"replica.slow_query_threshold",
}

func init() {
	rootCmd := &grumble.Command{
		Name: "app-env",
		Help: "app-env related commands",
		// TODO(wutao): print commonly-used app-envs
	}
	rootCmd.AddCommand(&grumble.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Help:    "list the app-envs binding to the table",
		Run: func(c *grumble.Context) error {
			return executor.ListAppEnvs(pegasusClient, useTable)
		},
	})
	rootCmd.AddCommand(&grumble.Command{
		Name: "set",
		Help: "set an env with key and value",
		Run: func(c *grumble.Context) error {
			if len(c.Args) != 2 {
				return fmt.Errorf("invalid number (%d) of arguments for `app-env set`", len(c.Args))
			}
			return executor.SetAppEnv(pegasusClient, useTable, c.Args[0], c.Args[1])
		},
		AllowArgs: true,
		Completer: func(prefix string, args []string) []string {
			/* fill with predefined app-envs */
			if len(args) == 0 {
				return filterStringWithPrefix(predefinedAppEnvKeys, prefix)
			}
			return []string{}
		},
	})
	rootCmd.AddCommand(&grumble.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Help:    "delete app-envs with specified key or key prefix",
		Run: func(c *grumble.Context) error {
			if len(c.Args) != 1 {
				return fmt.Errorf("invalid number (%d) of arguments for `app-env delete`", len(c.Args))
			}
			return executor.DelAppEnv(pegasusClient, useTable, c.Args[0], c.Flags.Bool("prefix"))
		},
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.BoolL("prefix", false, "to delete with key prefix")
		},
	})
	rootCmd.AddCommand(&grumble.Command{
		Name: "clear",
		Help: "clear all app-envs",
		Run: func(c *grumble.Context) error {
			return executor.ClearAppEnv(pegasusClient, useTable)
		},
	})
	shell.AddCommand(rootCmd)
}