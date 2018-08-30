inventory Management
--------------------

[![Build Status](https://travis-ci.org/Timothylock/inventory-management.svg?branch=master)](https://travis-ci.org/Timothylock/inventory-management)
[![codecov](https://codecov.io/gh/Timothylock/inventory-management/branch/master/graph/badge.svg)](https://codecov.io/gh/Timothylock/inventory-management)

An inventory management system made specifically for use of myself and used by other clubs. This is written in Go.

# Installing
TODO: Docker image

Required Go >= 1.8

## Environment Variables
Several environment variables must be set in order for this system to boot. They should all be self explanatory.
- DB_URL
- DB_USER
- DB_PASS
- DB_NAME

## Table Schema
A SQL script is included in [TODO](todo) which you must run to populate the database. I hope to incorporate this directly into the container in the future but that depends on the need to do so.

# API
Can be found [HERE](todo)

# Contributing
Want to contribute new features? Just open a pull request and I will be happy to look at it!

### CI
Travis is the one powering the continuous integration.

Writing tests is a _must_. No tests, no review. Existing tests must also pass before GitHub will allow merging.

While not enforced by Github (yet?), code coverage also must be reasonable.