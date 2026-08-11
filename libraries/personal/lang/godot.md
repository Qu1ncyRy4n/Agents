---
tldr: Respect Godot scene and asset ownership while validating the affected gameplay path.
---
# Godot

## Project Shape

Inspect the Godot version, `project.godot`, autoloads, scenes, scripts, addons,
and exported resources before changing behavior. Keep scene ownership clear:
edit the scene, script, or resource that owns the behavior rather than patching
around it elsewhere.

## GDScript

Follow the project's GDScript style for typed variables, signals, node paths,
and exported properties. Prefer explicit signal wiring and small scripts over
large scene controllers that accumulate unrelated behavior.

## Assets

Do not rename, move, reimport, or delete assets casually. Godot resource paths
are part of the project graph, and `.import` or generated metadata may change
when assets move.

## Validation

Run the project's documented Godot check, headless test, or editor smoke path
when available. For gameplay changes, verify the affected scene and describe
the manual path tested.

## Status

This is a prospective module. The current captured corpus did not include a
Godot repo, so these rules are placeholders for future dogfooding rather than
derived rules.
