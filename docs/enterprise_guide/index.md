---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://about.gitlab.com/handbook/engineering/ux/technical-writing/#assignments
---

# Enterprise guide for deploying and scaling a GitLab Runner Fleet

GitLab Runners, a software agent, executes the build jobs defined in your GitLab CI pipeline either using a local shell or container. Hosting a runner fleet requires a well-planned infrastructure that needs to include considerations for computing capacity, storage, network bandwidth and throughput. 

This guide provides a structured approach to developing a GitLab Runner deployment strategy based on your organization's requirements. The guide does not make specific recommendations regarding the type of infrastructure you should use for hosting a runner fleet. However, it provides tips and insights from our experience in operating the runner fleet on GitLab.com, which processes millions of CI jobs monthly.

## Planning Checklist

- [ ] Create a list of the number of teams that will be using GitLab CI.
- [ ] Catalog the programming languages, web frameworks, libraries in use at your organization. For example (GoLang, C++, PHP, Java, Python, JavaScript, React, Node.js).
- [ ] Estimate the number of CI jobs each team may execute per hour, per day.
- [ ] Validate if any team has build environment requirements that cannot be addressed using containers. 
- [ ] Validate if any team has build environment requirements that are best served by having runners dedicated to that team.
- [ ] Estimate the compute capacity that you may need to support the expected CI demand.

We expect that there will be different infrastructure stacks chosen to host runner fleets, (public cloud, on-premise). It is critical to note that the performance of the CI jobs on the runner fleet is directly related to the fleet's environment. Hosting the fleet on a shared computing platform is not recommended for organizations executing a high number of resource-intensive CI jobs.

## Understanding workers, executors and autoscaling capabilities.

Let's take a minute to discuss in more detail some key concepts about the runner. As described earlier, `gitlab-runner` is the exectuable that executes your ci jobs. That's true, but its also a very flexible and powerful piece of software. Each runner is an isolated process responsible for picking up requests for job executions and dealing with them according to pre-defined configurations. As an isolated process, each runner can create 'sub-processes' (also called machines or workers) to run jobs. In the next section we describe setting up a most basic runner configuration so as to set the stage for discussions regarding more advanced configuration options.

### Basic configuration (one runner, one worker)

-For the most basic configuration lets say you install the GitLab Runner software on a supported host and operating system. 

- After the install is complete, you execute the runner registration command just once and you select the `shell` executor. In the runner `config.toml` you set concurrency = 1.

```toml
concurrent = 1

[[runners]]
  name = "instance-level-runner-001"
  url = ""
  token = ""
  executor = "shell"

```

In this very basic configuration, the GitLab CI jobs that this runner can process will directly execute the host system on which you installed the runner. It's as if you were running the CI job commands yourself in a terminal. In this case, since you only executed the registration command one time, there will only be one [[runners]] section in `config.toml`. Assuming we now set the concurrency value to 1, this means that there will only be one runner  `worker` to execute CI jobs for the runner process on this system.

### Intermediate configuration (one runner, multiple workers)

- Taking it a step further, you can also register multiple runner `workers` on this same machine. 
- When this is done, you will notice multiple [[runner]] sections in the in the runner's `config.toml` file. 
- Assuming that all of these additional runner `workers` are registered to use the shell executor, and we update the value of the global configuration option, `concurrent` to 3, this means that the upper limit of jobs that can run concurrently on this host is equal to three. 

```toml
concurrent = 3

[[runners]]
  name = "instance_level_shell_001"
  url = ""
  token = ""
  executor = "shell"

[[runners]]
  name = "instance_level_shell_002"
  url = ""
  token = ""
  executor = "shell"

[[runners]]
  name = "instance_level_shell_003"
  url = ""
  token = ""
  executor = "shell"

```

So you can register many runner `workers` on the same machine as each one is an isolated process. Of course the performance of the CI jobs that each worker is completely dependent on the compute capacity of the host system. In the next section we describe more advanced configurations for autoscaling runners.

### Autoscaling configuration (one or more runner managers, multiple workers)

In an autoscaling setup, you can configure the runner to act only as a manager using the `docker+machine` or `Kubernetes` executor. In this type of manager only configuration, this runner agent is itself not executing any CI jobs. 

- **Docker+Machine executor:** in the `docker+machine` executor's case, the runner managers role is to, on-demand, provision VM instances with Docker installed. The CI jobs are executed on these virtual machines using Docker. We recommend testing the performance of your CI jobs on various machine types. You will also need to determine whether to optimize your choice of compute hosts based on speed or cost. 

- **Kubernetes executor:** in the Kubernetes executor's case, the runner managers role is to provision Pods on the target Kubernetes cluster. Each Pod is where the CI job is executed, and is comprised of multiple containers. The Pods used for job execution will typically require more compute and memory resources than the POD used for hosting the runner manager.

## Configure runners for your organization

### Instance Level - Shared Runners

- Starting with instance level (shared runners) in an autoscaling configuration is an efficient and effective option for providing a CI build infrastructure for your organization. However, from the analysis of your organization's requirements, you may determine that you need to plan for delivering runners at the group level. 

- The compute capacity of the infrastructure stack that you use for hosting the VM's or PODS that will execute the CI jobs will depend on the requirements you have captured as part of the planning checklist exercise and the technology stack you will use to host your runner fleet.

- It is also very likely that you will need to adjust the computing capacity once you actually start running CI workloads and analyze the performance data over time. On GitLab SaaS, we provision GCP n1-standard-1 instances with 3.75GB of RAM for CI jobs on Linux with the `docker executor`.

- For configurations using instance-level, Shared Runners with an autoscaling executor, we recommend that you configure at minimum two runner managers to start. On GitLab.com we have been able to scale to millions of jobs per month using five runner managers. A snippet of the `config.toml` configuration file for GitLab.com is provided [here](https://docs.gitlab.com/ee/user/gitlab_com/#configtoml). 

The total number of runner managers that you may need over time will depend on the following factors:

- The compute resources of the stack hosting the runner managers.
- The concurrency that you choose to configure for each runner manager.
- The load that is generated by the CI jobs that each manager is executing, hourly, daily, monthly.

## Monitoring GitLab Runners

An essential step in operating a GitLab Runner fleet at scale, as is our experience handling millions of CI jobs monthly on GitLab.com, is to set up and use the [**GitLab Runner monitoring**](../monitoring/index.md) capabilities included with GitLab.

{What we are missing is a step by step tutorial on how to monitor GitLab CI on the RAILS side and the GitLab Runner side.}

### How to prepare the Prometheus monitoring stack.

1. Step 1:Enable Prometheus on each runner manager.
1. Step 2:{placeholder}
1. Step 3:{placeholder}
1. Step 4:{placeholder}

### Alerting

The team at [Radio France](https://medium.com/radio-france-engineering/on-demand-ci-cd-with-gitlab-and-kubernetes-1d395105ac45) has configured the following alerts in their monitoring system. This is a resonable starting point for establishing a solid monitoring framework that you can use to effectively operate a GitLab Runner fleert at scale.

- Alert when the number of pending jobs is above X for more than Y minutes.
- Alert when the rate of GitLab “system” failure is above X for more than Y minutes. 

The logic behind the system failure metric is that, a high rate of `system failures` could indicate a more systemic issue with the runner fleet.

## Example Runner Fleet Configurations

The following section provides example customer implementations of runner fleets.

### Customer A: 

- Multiple GitLab instances, 5000+ projects
- Tech stack = OpenStack, VMWare, OpenShift.
- Instance-Level (Shared Runners): The customer's operations team provides Shared Runners for everyone. In one instance, there are ~60 Shared Runners deployed.
- Group Runners: At the group (team) level, each team manages their own runners.
- There is a dedicated Shared Runner cluster on OpenStack for hosting the instance-level (Shared Runners). 
    - Prometheus metrics
    - Basic shared runners offered with various virtual machine sizes:
      - small: 1 vCPU, 4GB memory, 20GB storage
      - medium: 2 vCPU, 4GB memory, 20GB storage 
      - large: 4 vCPU, 8GB memory, 40GB storage
      - large: 8 vCPU, 16GB memory, 50GB storage
 
