# Upto

Upto is a project based around the observation that in production systems with
many components, it can be very difficult to track and visualize a coordinated
many-step process, especially if the orchestration is distributed or difficult
to monitor. When you get asked (for example), why a server startup, build,
deployment, database downtime or other process is slower, it can be difficult to
track all the moving pieces.

Upto has the intention of solving this by providing a service that you can run and
report events to, either by UDP or HTTP, to provide a Gantt chart visualization and REST API.
As a secondary goal, it would be nice to support both push and pull mechanics,
allowing upto to gather information more easily over a variety of network setups.

The name is based on the question, "what is it up to?" and it is inspired by work
on two projects, one involving many tools, and the other being a build system that
required timings for sub-job level items.

 
