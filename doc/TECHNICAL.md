# Technical

## Outages

I've rolled bugs into production systems numerous times. These days serious bugs should show up quickly in KPIs through instrumentation. And progressive rollouts coupled with consistent, even automated rollbacks triggered by monitoring should mitigate most of the damage by these kinds of things.

Note I said _should_. As the complexity of system increases, safeguard systems can't always protect against unexpected interactions or prevent unexpected behavior. Once a system crosses the "deterministic machine" line into "living system" territory (when a live service becomes a live ecosystem), things get dicey.

Google, for example, has some of the most sophisticated infrastructure for progressive rollouts, canaries, etc. But it doesn't always mitigate damage and it requires payment of a significant price... velocity.

And game companies rushing to market against their competion can't really afford to sacrifice velocity. Taking it further, I'm tempted to assert that only speed matters. A really fast inner loop and release cadence contributes to a feeling of freshness, variety, and growth. It is, by definition, creative. Slowing the pace of development by imposing checkpoints and safeguards numbs the mind, making it hard for game developers to do their best work.

I've seen this pattern a few times... game company gets big and invests heavily in infrastructure and automation. It wants to leverage all that investment, all that prodigious infrastructure to launch a new title. Why would they want to build it again? And since they already have it... why not bake it in from the beginning?

But that wasn't the recipe that produced the hit in the first place. They recipe was (almost always) start the kernel of something new, small, and agile... improvise... and then grow.

Politics, the sunken cost fallacy, hubris, and so many other things can get in the way.

Automation works to prevent outages in the middle ground between living system and deterministic prototype. Up to that point, basic monitoring will help and teams should consider it a necessity. Some basic canary release ability should also become standard practice along the way. But up to that point and beyond... play the game. Use the software. Test it. Experience it. Games are not for robots. They are for humans. At least for now. Do ivory tower testing, but not exclusively. Nothing can replace a developer with sufficient time and attention to stay in touch with the machine.

## Containers

What's the point? Change control. Consistency. Security. Microservicability (my word). Scalability.

Containers separate machine management, configuration, composition, security, and scalability from the application developer and offloads that to the infrastructure provider.

Scaling requires services to talk to each other over the network. At cloud scale, monoliths don't work. And at cloud scale, you can't expect a homogenous fleet. And you can't expect all vendors to develop for the same base operating system.

Google Borg (the progenitor of Kubernetes) ingests an increasingly varied fleet of hardware. That hardware gets pooled and exported as fungible compute capacity. That's what app developers want... capacity. They generally don't care about the internals. If they do, they probably shouldn't.

Game developers are notorious holdouts in this regard, wanting to stay close to metal for maximum raw performance and low latency. But more and more game developers will move to public cloud because of the cost of rapidly scaling up (and down) to service the needs of an industry with equally notoriously short cycle times.

Nobody wants a billion dollar data center sitting idle. Game companies can learn to sell that excess capacity or avoid acquiring it in the first place. And modern containerized CGroup virtualization doesn't really cost that much.

Real-time, low-latency messaging my still require unique soltuions, but most everything else can live in the Cloud.

## Availability

A system is either "available" or it "isn't" given a threshold definition of availability. If a service _should_ be able to serve 100k concurrent users, but can't (it can only serve 50k) for a period of time, SRE / ops considers that service not available at the moment of the sample. It's an instantaneous metric and defined using "indicators" or "signals".

CCUs aren't necessarily the best indicator of availability. A problem with a login service might limit the CCUs into a service which then responds to most requests successfuly.

A service-level objective (SLO) would define the sample rate for the indicator, the time-slices that matter, and the hit/miss ratio allowed for the time slice. The rolling average of missed slices is what matters for the SLO. When the rolling average dips below an alerting threshold, the service is "out of slo" and someone needs to take action.

But no service should target 100% availability. That's not reasonably possible or realistically desirable. Extreme high availability can indicate that something is out of balance, or measured or defined incorrectly.

For internal services, extreme HA can lead consumers to a false sense of security. When the service goes down now and then, consumers will have to add resilience to their services.

## Future Tech for Me / Hobby Projects

LLMs, ML for CV, Robotics. I have a pile of O'Reilly books on MLOps and Techniques waiting for me to dig into. I have built and operate Raspberry Pi IoT cloud security cameras for my home that support small ML models for CV. I have a few microcontroller projects that I improve now and then. I have a Mecanum Car Chassis waiting for me to control with a Raspberry Pi.

I do lots of home improvement projects, with light woodworking.

## Self-Assessment

| Tech           | Rating |
| -------------- | ------ |
| Linux          | 6      |
| Python         | 6      |
| Go             | 7      |
| Docker         | 5      |
| Open Telemetry | 5      |
| AWS ECS        | 5      |
| AWS Serverless | 2      |
| Terraform      | 2      |
| git            | 7      |
| GitLab         | 2      |
