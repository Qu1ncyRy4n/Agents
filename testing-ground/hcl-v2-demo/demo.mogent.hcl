# Runnable local v2 dogfood fixture. The exploratory mogent.hcl beside this
# file intentionally retains design notes and is not a valid configuration.

mogent {
  format = 2
}

sources {
  source "cdint" {
    local = "../../../cdint-demo-lib"
  }
}

outputs {
  output "agent-instructions" {
    paths = ["AGENTS.md", "CLAUDE.md"]
    kind  = "markdown"

    source "cdint" {
      from = "cdint:agents"

      select = {
        org = {
          baseline = true

          core = {
            accept_defaults = true
          }

          optional = {
            "project-defaults" = true
          }
        }

        else = "exclude"
      }
    }
  }
}
