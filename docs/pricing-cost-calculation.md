---
title: Pricing & Cost Calculation
description: How traceknot turns token usage into a dollar cost, and what to do if a number looks wrong.
---

# Pricing & Cost Calculation

## How is cost calculated?

traceknot matches each step's model and token usage against a maintained
list of current provider prices, and adds it up per session. That covers
the input/output split, the discounted rate for cached tokens, and any
extra per-use fees a provider charges for specific tools (for example, some
providers bill web search separately from regular tokens).

## Why does a session show no cost, or $0.00?

That means the model used in that session isn't in the pricing list yet —
almost always because it's a model that launched very recently. New models
are typically added within a day of release. Once the update lands, cost
for that model is calculated correctly going forward, and **sessions
already recorded with that model are recalculated too** — you don't need to
re-run anything.

## If a number doesn't add up

Pricing lists are maintained by hand and can lag a provider's own pricing
page, or miss a discount tier. If a session's cost looks wrong to you,
[open an issue on GitHub](https://github.com/joroden/traceknot/issues) with
the session id and what you expected — that's the fastest way to get a
pricing fix shipped.
