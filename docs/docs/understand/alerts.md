---
title: Alerts
sidebar_position: 60
---

# Alerts from Minder

Alerts are a core feature of Minder providing you with notifications about the status of your registered
repositories. These alerts automatically open and close based on the evaluation of the rules defined in your profiles.

When a rule fails, Minder opens an alert to bring your attention to the non-compliance issue. Conversely, when the
rule evaluation passes, Minder will automatically close any previously opened alerts related to that rule.

In the alert, you'll be able to see details such as:
* The repository that is affected
* The rule type that failed
* The profile that the rule belongs to
* Guidance on how to remediate and also fix the issue
* Severity of the issue. The severity of the alert is based on what is set in the rule type definition.

## Alert types

Minder supports alerts of type GitHub Security Advisory.

The following is an example of how the alert definition looks like for a give rule type:

```yaml
---
version: v1
type: rule-type
name: artifact_signature
...
def:
  # Defines the configuration for alerting on the rule
  alert:
    type: security_advisory
    security_advisory:
      severity: "medium"
```

## Configuring alerts in profiles

Alerts are configured in the `alert` section of the profile yaml file. The following example shows how to configure
alerts for a profile:

    ```yaml
    ---
    version: v1
    type: profile
    name: github-profile
    context:
      provider: github
    alert: "on"
    repository:
      - type: secret_scanning
        def:
          enabled: true
    ```

The `alert` section can be configured with the following values: `on`, `off` and `dry_run`. The default value is `on`.
