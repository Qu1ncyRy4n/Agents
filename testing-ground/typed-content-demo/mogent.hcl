mogent { format = 2 }

sources {
  source "demo" { local = "library" }
}

outputs {
  output "package" {
    path = "agent-export"
    kind = "dir-tree"

    source "demo" {
      from      = "demo:physical"
      node      = "review"
      operation = "copy"
      into      = ".agents/skills/review"
    }

    source "demo" {
      from      = "demo:physical"
      node      = "review/SKILL.md"
      heading   = ["Review", "Procedure"]
      operation = "render-markdown"
      into      = "AGENTS.md"
    }

    source "demo" {
      from      = "demo:physical"
      node      = "settings.toml"
      operation = "copy"
      into      = "config/tool.toml"
    }
  }
}
