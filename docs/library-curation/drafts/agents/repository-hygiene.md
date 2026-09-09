# Keep Local State And Generated Binaries Out Of Git

Repeated verbatim in wire-lab, current cdint-grid, and promisegrid source guides:

> Do not commit local state files (for example `.grok`, `.grok.lock`) or
> generated binaries.

Source: [`promisegrid_wire-lab_refs_heads_main_AGENTS.md`](../reference/source-guides/promisegrid_wire-lab_refs_heads_main_AGENTS.md#project-structure--module-organization)

Secrets and credentials will be considered when the later source section that
contains them is reviewed. A future cleanup skill may supplement this
always-loaded prohibition; it must not replace it.
