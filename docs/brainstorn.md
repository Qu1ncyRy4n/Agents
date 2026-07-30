# brainstorm

## YAML-based 

### Processing

1. Start at the root of the repo and look for a `.mogent/config.yaml`
   file.
2. If found, read the config file and process any `include`
   directives.  An included file can be a local file path or a URL.
   Included files can themselves include other files.
3. If any definitions conflict, the last definition wins.
4. After processing all includes, the final definitions in the config
   file override any definitions in the included files.
5. After processing the config file, the tool will look for any `.md`
   files referenced in the compiled config and concatenate them into a
   single output file, in the order specified by the config.
6. Before writing the output file, the tool will process any Go templates in the
   .md files, substituting any variables defined in the config file or
   provided by the tool.
7. Write the final output to a file named `AGENTS.md` in the root of the repo.


### {repo}/.mogent/config.yaml:

All paths are relative to the same directory as the config file.

```
config:
    - include: http://github.com/ciwg/agents/main/ciwg.yaml
    - include: ~/.mogent.yaml

# following definitions are repo-specific and override any definitions in the included files
category:
    format:
        module:
            name: format
            source: repo-specific-format.md  # this path is {repo}/.mogent/repo-specific-format.md
    wire-lab:
        module:
            - name: sims
              source: sims.md  # this path is {repo}/.mogent/sims.md
            - name: pocs
              source: pocs.md  # this path is {repo}/.mogent/pocs.md
```

### templating

- Each included .md file is treated as a Go template.
- Dev can define variables in the config file. For example, if the
  config file defines a variable `project_name`, you can use 
  `{{ project_name }}` in your .md files to insert its value.
- Tool can define default variables, such as `repo_name`, `repo_url`,
  etc., that can be used in the templates. 

## Module Category Model

Current compact render order:

1. Identity
   - agent role and project context
   - project overview
   - tech stack
   - project structure
2. Instructions
   - workflow
   - code changes
   - testing
   - commits
   - decision protocol
   - thought experiments
   - DR/DI protocol
   - comment preservation
   - TODO tracking
3. Constraints
   - never-do rules
   - always-do rules
   - security and prohibited actions
   - runtime artifact hygiene
   - hard compliance rules
4. Format
   - coding style
   - diff discipline
   - error handling format
   - response and handoff format
   - glossary

Richer module-library categories to explore:

- Cognition / process intent: thinking depth, TE/DI/DF behavior, fast iteration vs deliberate research, learning-focused mode, and "think like a..." rules.
- Communication / chat style: directness, Socratic mode, TTS-friendly output, whiteboard workflow, humor, and simplicity level.
- Code: language-specific style, stack rules, testing strategy, commit cadence, security, and developer involvement level.
- Notes / docs: README conventions, changelog style, dev logs, Obsidian/session notes, human-facing docs vs LLM-only docs, and TODO/DR/DI conventions.

The current generated `AGENTS.md` uses the compact four-section structure. The richer categories should become optional modules, presets, or subtrees once TUI selection and saving are usable.

## TUI Selection Model

Selection without saving is useful as a preview/staging state:

- try combinations without changing repo config,
- preview what would be included,
- back out without touching `AGENTS.toml`,
- later compare the temporary selection against the saved selection.

For normal dogfooding, selection needs save support soon. The intended flow is:

1. Toggle blocks in `mogent tui`.
2. Show dirty state when the in-memory selection differs from `AGENTS.toml`.
3. Preview rendered output.
4. Press `s` to persist selected blocks back to `AGENTS.toml`.
5. Run or offer `mogent build` to regenerate `AGENTS.md`.

## Fork Import Workflow

Useful changes can be imported from someone else's fork without opening a pull request.

Inspect full branch changes:

```sh
git remote add theirname https://github.com/theirname/repo.git
git fetch theirname
git checkout -b import-their-changes
git merge theirname/branch-name
git diff main...HEAD
```

Import only specific commits:

```sh
git fetch theirname
git cherry-pick <commit-sha>
```

Import one file:

```sh
git fetch theirname
git checkout theirname/branch-name -- path/to/file
```

Before committing imported work, review license/authorship, check for secrets or local machine paths, and test on an import branch.
