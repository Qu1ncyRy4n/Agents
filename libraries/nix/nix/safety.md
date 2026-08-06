---
tags: [nix/safety, risk/system-activation, scope/project]
tldr: Do not activate or touch host-critical Nix settings without explicit scope.
priority: 1.0
scope: project
---
# Safety

Do not run an activating system command without explicit user approval. Do not
change hardware configuration, bootloader, storage, firewall, VPN, DNS,
networking, shell startup behavior, or secrets unless the request explicitly
covers that area.
