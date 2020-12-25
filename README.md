# cliconfig
Combines Viper and Cobra libraries for flexible configurable values.

```go
  &cobra.Command{
    PreRunE: func(cmd *cobra.Command, args []string) error {
      return cliconfig.BindViperDefaults(cmd, "git")
    },
    RunE: func(cmd *cobra.Command, args []string) error {
    log.Info().Msg("Submitting git modification")
    gmc := gitModifyConfig{}
                          if err := cliconfig.Populate(cmd, &gmc); err != nil {
                                  return err
                          }
                          fmt.Printf("%+v", gmc)
                          return nil
                  },
          }
```
