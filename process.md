## process workflow
```mermaid
---
title: process workflow
---
flowchart LR
	init-->|"start"|started
	stopped-->|"start"|started
	started-->|"cancel"|cancelled
	stopped-->|"cancel"|cancelled
	started-->|"complete"|completed
	started-->|"stop"|stopped
```