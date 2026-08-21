# Explore the Whole Decision Surface Before Locking Local Syntax

Status: proposed module; invented during v2 protolibrary design.

Before implementing a provisional schema or command, inventory the neighboring
decisions that could make it wrong. Capture a small candidate vocabulary,
several real examples, safety boundaries, migration consequences, and open
questions first. Then narrow the field with a source-grounded thought
experiment and record the owner decision.

This is breadth-first decision work: it prevents the first convenient syntax
from becoming an accidental permanent architecture. It does not require fully
solving every adjacent question before useful extraction, prototyping, or
evidence gathering continues.

Proposed example: record inline imports, manifest-only composition, choice
rules, directory inheritance, cycle safety, and migration together; postpone
implementation; continue extracting modules until real cases reveal which
parts are necessary.
