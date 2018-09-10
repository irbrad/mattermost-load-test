package main

import (
	"os"
	"path/filepath"

	"github.com/mattermost/mattermost-load-test/kubernetes"
	"github.com/mattermost/mattermost-load-test/ltops"
	"github.com/mattermost/mattermost-load-test/terraform"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var createCluster = &cobra.Command{
	Use:   "create",
	Short: "Creates a cluster to run Mattermost on for load testing",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}

		return nil
	},
	RunE: createClusterCmd,
}

func createClusterCmd(cmd *cobra.Command, args []string) error {
	var config ltops.ClusterConfig
	config.Name, _ = cmd.Flags().GetString("cluster")
	config.AppInstanceType, _ = cmd.Flags().GetString("app-type")
	config.AppInstanceCount, _ = cmd.Flags().GetInt("app-count")
	config.DBInstanceType, _ = cmd.Flags().GetString("db-type")
	config.DBInstanceCount, _ = cmd.Flags().GetInt("db-count")
	config.LoadtestInstanceCount, _ = cmd.Flags().GetInt("loadtest-count")
	clusterType, _ := cmd.Flags().GetString("type")

	workingDir, err := defaultWorkingDirectory()
	if err != nil {
		return err
	}
	config.WorkingDirectory = filepath.Join(workingDir, config.Name)

	if clusterType == kubernetes.CLUSTER_TYPE {
		_, err = kubernetes.CreateCluster(&config)
		if err != nil {
			return err
		}
	} else if clusterType == terraform.CLUSTER_TYPE {
		if len(config.AppInstanceType) == 0 {
			return errors.New("required flag \"app-type\" not set")
		}
		if len(config.DBInstanceType) == 0 {
			return errors.New("required flag \"db-type\" not set")
		}

		_, err = terraform.CreateCluster(&config)
		if err != nil {
			return err
		}
	} else {
		return errors.New("unrecognized type: " + clusterType)
	}

	return nil
}

func init() {
	createCluster.Flags().StringP("cluster", "c", "", "a unique name for the cluster (required)")
	createCluster.MarkFlagRequired("cluster")

	createCluster.Flags().StringP("type", "t", "", "the type of cluster, terraform or kubernetes (required)")
	createCluster.MarkFlagRequired("type")

	createCluster.Flags().String("app-type", "", "the app instance type (required for terraform)")

	createCluster.Flags().Int("app-count", 1, "the number of app instances")

	createCluster.Flags().String("db-type", "", "the db instance type (required for terraform)")

	createCluster.Flags().Int("db-count", 1, "the number of db instances")

	createCluster.Flags().Int("loadtest-count", 1, "the number of loadtest instances")

	createCluster.Flags().SortFlags = false

	rootCmd.AddCommand(createCluster)
}
