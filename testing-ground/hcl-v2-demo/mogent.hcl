# HCL v2 design spike. Current Mogent does not execute this file.
# This retains the hierarchy of the current YAML manifest: sources, then
# independent outputs, then ordered selection for each output.

mogent {
  format = 2
}

sources {
  source "cdint" {
    local = "../../../cdint-demo-lib"
  }

  # This is a real, immutable QMR revision. A later lockfile design may record
  # the resolved commit outside this authored configuration.
  source "qmr" {
    git    = "https://github.com/Qu1ncyRy4n/qmr-agents-library.git"
    commit = "14724d4191c9e026bc5a1e8549b7a96cc0c970e1"
  }

}

outputs {
  output "agents" {
    path = "AGENTS.md"
    kind = "markdown"

    # Entries render in this order. The organization library cannot prescribe
    # selections for this repository; the consumer owns this list.
    include {
      all = "cdint:org/mandatory"
    }

    include {
      source = "cdint:org/review"
    }

    include {
      all = "cdint:project/all"
    }

    # Metadata selection is declarative and scoped to the selected source.
    # `all` requires each tag; `any` matches one or more listed tags.
    include {
      tags {
        source = "qmr:agents"
        all    = ["language/go", "workflow"]
      }
    }

    include {
      all = "qmr:agents/accessibility"
    }

    include {
      all = "qmr:agents/style"
    }
  }

  output "skills" {
    path = ".agents/skills"
    kind = "tree"

    include {
      all = "qmr:skills"
    }
  }
}



### Heres a different sketch, what do yo uthink? I guess this is more nix like or something

# HCL v2 design spike. Current Mogent does not execute this file.
# This retains the hierarchy of the current YAML manifest: sources, then
# independent outputs, then ordered selection for each output.

mogent {
  format = 2
}

sources {
  source "cdint" {
    local = "../../../cdint-demo-lib"
  }

  # This is a real, immutable QMR revision. A later lockfile design may record
  # the resolved commit outside this authored configuration.
  source "qmr" {
    git    = "https://github.com/Qu1ncyRy4n/qmr-agents-library.git"
    commit = "14724d4191c9e026bc5a1e8549b7a96cc0c970e1"
  }

  source "local" {
    local = "."
  }
}

outputs {
  output "agents" {
    path = "AGENTS.md"
    kind = "markdown"

    # Entries render in this order. The organization library cannot prescribe
    # selections for this repository; the consumer owns this list.
    include {
      all = "cdint:org/mandatory"
      source = "cdint:org/review"
      all = "cdint:project/all"
      tags {
        source = "qmr:agents"
        all    = ["language/go", "workflow"]
      }
      all = "qmr:agents/accessibility"
      all = "qmr:agents/style"
    }
  }


  output "skills" {
    path = ".agents/skills"
    kind = "tree"
    include {
      all = "qmr:skills" {
        except tags ["os/nix" "lang/nix"]
      }
      all = "cdint:skills"


    }
  }
}


# This doesn't really look quite right. Here is what I'd want, maybe from a different langage
# ideal would be somethign like:

outputs {
  agents = {
    path = "AGENTS.md"
    aliases = [ "CLAUDE.md" "GEMINI.md"]
    kind = "markdown"

    include {
      from cdint:agents {
        include: org/madantory, project/ # project/all
        exclude section heading2/subheading3
      }
      from qmr:agents {
        include:
          tags: ["langauge/go"., "workflow"]
          mods: ["os/nix" "lang/nix"] # not sure what best name here would be
      }

    }
  skills = {
    include {
      qmr:skills = all
      cdint:skills = all
    }
  }

  }
}
