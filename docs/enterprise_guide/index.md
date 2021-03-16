---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://about.gitlab.com/handbook/engineering/ux/technical-writing/#assignments
---

# Considerations for scaling a fleet of runners for your enterprise

When you host a fleet of runners, you need a well-planned infrastructure that takes
into consideration your:

- Computing capacity.
- Storage capacity.
- Network bandwidth and throughput. 

This guide provides a structured approach to developing a GitLab Runner deployment strategy
based on your organization's requirements. The guide does not make specific recommendations
about the type of infrastructure you should use. However, it provides tips and insights from
the experience of operating the runner fleet on GitLab.com, which processes millions of
CI/CD jobs each month.

## Planning Checklist

- <input type="checkbox"> Create a list of the teams that will use GitLab CI/CD.
- <input type="checkbox"> Catalog the programming languages, web frameworks, and libraries in use
  at your organization. For example, GoLang, C++, PHP, Java, Python, JavaScript, React, Node.js.
- <input type="checkbox"> Estimate the number of CI/CD jobs each team may execute per hour, per day.
- <input type="checkbox"> Validate if any team has build environment requirements that cannot be
  addressed by using containers. 
- <input type="checkbox"> Validate if any team has build environment requirements that are best served
  by having runners dedicated to that team.
- <input type="checkbox"> Estimate the compute capacity that you may need to support the expected demand.

You will probably need different infrastructure stacks to host different runner fleets.
For example, you may want some runners in the public cloud and some on-premise.

The performance of the CI/CD jobs on the runner fleet is directly related to the fleet's environment.
If you are executing a large number of resource-intensive CI/CD jobs, hosting the fleet on a shared
computing platform is not recommended.

## Workers, executors, and autoscaling capabilities

The `gitlab-runner` executable is a flexible and powerful piece of software that runs your CI/CD jobs.
Each runner is an isolated process responsible for picking up requests for job executions and dealing
with them according to pre-defined configurations. As an isolated process, each runner can create
"sub-processes" (also called "machines" or "workers") to run jobs.

The next section describes how to set up a runner with very few configuration options.
You can start with this configuration and then build on it.

### Basic configuration: one runner, one worker

For the most basic configuration, you install the GitLab Runner software on a supported host and operating system. 

After the installation is complete, you execute the runner registration command just once
and you select the `shell` executor. Then you edit the runner `config.toml` file to set concurrency to `1`.

```toml
concurrent = 1

[[runners]]
  name = "instance-level-runner-001"
  url = ""
  token = ""
  executor = "shell"
```

The GitLab CI/CD jobs that this runner can process are executed directly on the host system where you installed the runner.
It's as if you were running the CI/CD job commands yourself in a terminal. In this case, because you only executed the registration
command one time, the `config.toml` file contains only one `[[runners]]` section. Assuming we set the concurrency value to `1`,
only one runner "worker" can execute CI/CD jobs for the runner process on this system.

### Intermediate configuration: one runner, multiple workers

Taking it a step further, you can also register multiple runner workers on this same machine. 
When this is done, the runner's `config.toml` file has multiple `[[runner]]` sections in it. 
If all of the additional runner workers are registered to use the shell executor,
and we update the value of the global configuration option, `concurrent`, to `3`, 
the upper limit of jobs that can run concurrently on this host is equal to three. 

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

You can register many runner workers on the same machine, and each one is an isolated process.
The performance of the CI/CD jobs for each worker is dependent on the compute capacity of the host system.

In the next section we describe more advanced configurations for autoscaling runners.

### Autoscaling configuration: one or more runner managers, multiple workers

In an autoscaling setup, you can configure a runner to act as a manager of other runners.
You can do this with the `docker-machine` or `Kubernetes` executor only. In this type of
manager-only configuration, the runner agent is itself not executing any CI/CD jobs. 

- **Docker Machine executor:** The runner manager provisions on-demand virtual machine instances that have Docker installed.
  On these VMs, Docker executes the CI/CD jobs. You should test the performance of your CI/CD jobs
  on various machine types. You should also decide whether to optimize your compute hosts based on speed or cost. 

- **Kubernetes executor:** The runner manager provisions pods on the target Kubernetes cluster.
  The CI/CD jobs are executed on each pod, which is comprised of multiple containers. The pods used for job execution
  typically require more compute and memory resources than the pod that hosts the runner manager.

## Configure runners for your organization

### Instance-level shared runners

Using instance-level shared runners in an autoscaling configuration is an efficient and effective way to start.
However, from the analysis of your organization's requirements, you may find that you need to create runners at the group level.

The compute capacity of the infrastructure stack where you host your VMs or pods depends on:

- The requirements you captured as part of the planning checklist exercise.
- The technology stack you use to host your runner fleet.

It is also likely that you will need to adjust the computing capacity after you start running CI/CD workloads
and analyzing the performance over time. On GitLab.com, we provision GCP n1-standard-1 instances with 3.75GB
of RAM for CI/CD jobs on Linux with the `docker` executor.

For configurations that use instance-level shared runners with an autoscaling executor,
we recommend that you start with, at minimum, two runner managers. On GitLab.com, we have been able to scale to
millions of jobs per month by using five runner managers. You can view
[a snippet of the `config.toml` configuration file for GitLab.com](https://docs.gitlab.com/ee/user/gitlab_com/#configtoml). 

The total number of runner managers that you may need over time depends on the following factors:

- The compute resources of the stack hosting the runner managers.
- The concurrency that you choose to configure for each runner manager.
- The load that is generated by the CI/CD jobs that each manager is executing hourly, daily, and monthly.

## Monitoring runners

An essential step in operating a runner fleet at scale is to set up and use the [runner monitoring](../monitoring/README.md) capabilities included with GitLab. 

The following table includes a summary list of GitLab Runner metrics. The list does not include the GoLang specific process metrics. To view those metrics on a runner, execute the command as noted [here](https://docs.gitlab.com/runner/monitoring/README.html#available-metrics).

| metric_name | description |
| ------ | ------ |
| gitlab_runner_api_request_statuses_total |The total number of API requests, partitioned by runner, endpoint and status. |
| gitlab_runner_autoscaling_machine_creation_duration_seconds  | Histogram of machine creation time.|
| gitlab_runner_autoscaling_machine_states  | The current number of machines per state in this provider. |
| gitlab_runner_concurrent | The current value of concurrent setting. |
| gitlab_runner_errors_total| The number of caught errors. |
| gitlab_runner_job_duration_seconds | Histogram of job durations. |
| gitlab_runner_jobs_total| cell |
| gitlab_runner_limit| The current value of concurrent setting. |
| gitlab_runner_request_concurrency | The current number of concurrent requests for a new job. |
| gitlab_runner_request_concurrency_exceeded_total| Counter tracking exceeding of request concurrency. |
| gitlab_runner_version_info| A metric with a constant '1' value labeled by different build stats fields.|
| process_cpu_seconds_total | Total user and system CPU time spent in seconds. |
| process_max_fds | Maximum number of open file descriptors. |
| process_open_fds| Number of open file descriptors. |
| process_resident_memory_bytes | Resident memory size in bytes. |
| process_start_time_seconds| Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes| Virtual memory size in bytes. |
| process_virtual_memory_max_bytes| Maximum amount of virtual memory available in bytes.|

### How to prepare the Prometheus monitoring stack

In this section we provide a step by step guide to configuring and using the Prometheus monitoring stack. For reference,  we provide the source code to some our [Grafana Dashboards](https://gitlab.com/gitlab-com/runbooks/tree/master/dashboards).

{What we are missing is a step by step tutorial on how to monitor GitLab CI on the RAILS side and the GitLab Runner side.}

1. Step 1:Enable Prometheus on each runner manager.
1. Step 2:{placeholder}
1. Step 3:{placeholder}
1. Step 4:{placeholder}

### Alerting

The team at [Radio France](https://medium.com/radio-france-engineering/on-demand-ci-cd-with-gitlab-and-kubernetes-1d395105ac45) has configured the following alerts in their monitoring system. This is a reasonable starting point for establishing a solid monitoring framework that you can use to effectively operate a runner fleet at scale.

- Alert when the number of pending jobs is above `x` for more than `y` minutes.
- Alert when the rate of GitLab “system” failure is above `x` for more than `y` minutes. 

The logic behind the system failure metric is that a high rate of `system failures` might indicate a systemic issue with the runner fleet.

## Complex runner deployment scenario (GitLab.com)

The following section summarizes a few key points regarding the runner architecture on GitLab.com.

GitLab.com uses multiple runner managers for each build environment (Linux+Docker, Windows). 
These multiple runners provides redundancy. If you have only one runner manager, then it is a single point of failure.
For additional fault tolerance, you can choose to host managers on different on-premise infrastructure stacks
or take advantage of multi-region capabilities offered by your public cloud provider.

GitLab.com has multiple runner managers because it provides runners with different characteristics.
The GitLab.com shared Linux runner managers are:

- Hosts with the `gitlab-runner` executable configured to autoscale using Docker Machine and the Google Compute API.
- The only virtual machines that are always active. The autoscaled virtual machines are one-time use only.
  They are used for only one CI/CD job and deleted immediately after the job completes.

The GitLab.com Prometheus environment ingests metrics from each runner manager. The GitLab.com infrastructure team
uses these metrics to monitor CI/CD job queues, and to operate and optimize the platform.

## Customer examples of runner fleet configurations

The following section provides example customer implementations of runner fleets.

### Customer A

- Multiple GitLab instances, 5000+ projects
- Tech stack = OpenStack, VMWare, OpenShift.
- Instance-level shared runners): The customer's operations team provides shared runners for everyone.
  In one instance, there are approximately 60 shared runners deployed.
- Group runners: At the group (team) level, each team manages their own runners.
- There is a dedicated shared runner cluster on OpenStack for hosting the instance-level shared runners. 
- Basic shared runners offered with various virtual machine sizes:
  - small: 1 vCPU, 4GB memory, 20GB storage
  - medium: 2 vCPU, 4GB memory, 20GB storage 
  - large: 4 vCPU, 8GB memory, 40GB storage
  - large: 8 vCPU, 16GB memory, 50GB storage
