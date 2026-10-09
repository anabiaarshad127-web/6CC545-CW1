# Implementation log

Record genuine work as it happens. Do not prefill outcomes.

| Date | Task / commit | Decision and reason | Test run and actual result | Feedback / next step | Sources or AI assistance |
|---|---|---|---|---|---|
| | | | | | |

| 2026-10-09 | Week 2 — `4392a6f` — Implement REST status resource | Added GET and POST at `/v1/status`. Kept request decoding and status state separate. Limited JSON bodies to 4096 bytes and rejected invalid input before changing state. | All 13 `TestWeek02_REST` cases passed. `go test ./... -count=1 -timeout 20s` passed. `go vet ./...` completed without reported problems. Live requests confirmed valid updates and rejection of an invalid status. | Review the next week's requirements before adding more communication features. | 6CC545 CW1 Assessment Specification and module guidance. ChatGPT assisted with explanations and troubleshooting; commands and results were checked by me. |
