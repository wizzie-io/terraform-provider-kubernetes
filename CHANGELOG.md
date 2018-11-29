# Change Log

## [v1.3.0-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.3.0-custom) (2018-11-23)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.2.2-custom...v1.3.0-custom)

**Fixed bugs:**

- TF kubernetes provider tries to in-place update not modifiable attribute [\#74](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/74)

**Closed issues:**

- Importing a config\_map volume with no mode causes panic [\#81](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/81)
- Importing a resource with envFrom: configMapRef crashes [\#76](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/76)
- some annotations are being re-applied on every terraform run [\#72](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/72)

**Merged pull requests:**

- EmptyDir ‘size\_limit’ causing quantity parse error [\#88](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/88) ([sl1pm4t](https://github.com/sl1pm4t))
- Fix TF re-adding annotations on every apply [\#87](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/87) ([sl1pm4t](https://github.com/sl1pm4t))
- Fix: StatefulSet pod\_management\_policy [\#86](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/86) ([sl1pm4t](https://github.com/sl1pm4t))
- Add/Fix attributes to/on service resource [\#85](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/85) ([sl1pm4t](https://github.com/sl1pm4t))
- Added ServiceExternalTrafficPolicy [\#84](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/84) ([sebastianroesch](https://github.com/sebastianroesch))
- Fixing panic when importing resource with config\_map volume [\#82](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/82) ([itmecho](https://github.com/itmecho))
- Add kubernetes\_secret datasoucre [\#80](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/80) ([Phylu](https://github.com/Phylu))
- Add missing empty\_dir size\_limit  [\#79](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/79) ([itmecho](https://github.com/itmecho))
- Fixing panic when loading pods with env\_from config\_map\_ref [\#77](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/77) ([itmecho](https://github.com/itmecho))
- Adds Deployment Data Source [\#75](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/75) ([itmecho](https://github.com/itmecho))
- fix affinity.pod\_\(anti\)\_affinity namespaces missing in flattener [\#71](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/71) ([kolach](https://github.com/kolach))

## [v1.2.2-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.2.2-custom) (2018-10-17)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.2.1-custom...v1.2.2-custom)

**Closed issues:**

- Crashing Terraform [\#65](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/65)
- Can't use variables/datasources for labels? [\#63](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/63)
- Pod dns\_config not recognised [\#62](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/62)
- deployment strategy doesn't seem to work as expected [\#61](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/61)
- Deployment pod affinity not implemented [\#44](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/44)
- Use affinity and tolerations with DaemonSet [\#26](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/26)

**Merged pull requests:**

- Fix “expected pointer, but got nil” error when Deployment strategy = Recreate [\#69](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/69) ([sl1pm4t](https://github.com/sl1pm4t))
- affinity, yet another attempt [\#68](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/68) ([kolach](https://github.com/kolach))
- Fix panic caused by DNSConfig attribute [\#66](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/66) ([sl1pm4t](https://github.com/sl1pm4t))
- Add `dns\_config` attribute to Pod Spec [\#64](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/64) ([sl1pm4t](https://github.com/sl1pm4t))
- Don't ForceNew for label changes. [\#58](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/58) ([anuraaga](https://github.com/anuraaga))

## [v1.2.1-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.2.1-custom) (2018-10-01)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.2.0-custom...v1.2.1-custom)

**Closed issues:**

- Feature Request: Support mountOptions in PersistentVolumes [\#54](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/54)
- kubeconfig created by a Terraform resource [\#47](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/47)

**Merged pull requests:**

- Add toleration support to pod spec [\#59](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/59) ([sl1pm4t](https://github.com/sl1pm4t))
- Add mount\_options attribute to PersistentVolume [\#56](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/56) ([sl1pm4t](https://github.com/sl1pm4t))
- Adding missing resources to README. [\#55](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/55) ([coryodaniel](https://github.com/coryodaniel))
- Added simpler build alternative [\#53](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/53) ([bennettellis](https://github.com/bennettellis))
- Static compilation for binaries [\#51](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/51) ([ThatsMrTalbot](https://github.com/ThatsMrTalbot))
- rbac api\_groups should be optional to allow for of non\_resource\_urls [\#42](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/42) ([abruehl](https://github.com/abruehl))

## [v1.2.0-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.2.0-custom) (2018-08-06)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.1.2-custom...v1.2.0-custom)

**Closed issues:**

- Support ClusterRole, Role, RoleBinding ClusterRoleBinding [\#35](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/35)
- Support k8s role and rolebinding [\#9](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/9)

**Merged pull requests:**

- Add RBAC resources [\#39](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/39) ([sl1pm4t](https://github.com/sl1pm4t))
- Set volume of configmap default mode value. [\#38](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/38) ([breeze7086](https://github.com/breeze7086))

## [v1.1.2-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.1.2-custom) (2018-07-18)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.1.1-custom...v1.1.2-custom)

**Closed issues:**

- jsonpatch error on ingress updates [\#31](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/31)
- 1.1.0 Linux Version not working [\#25](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/25)

**Merged pull requests:**

- Implement Update for `cron\_job` resources. [\#34](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/34) ([sl1pm4t](https://github.com/sl1pm4t))
- Cronjob / Job backoff limit attribute [\#33](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/33) ([sl1pm4t](https://github.com/sl1pm4t))
- CronJob: Fix starting\_deadline\_seconds successful\_job\_history\_limit [\#32](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/32) ([sl1pm4t](https://github.com/sl1pm4t))

## [v1.1.1-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.1.1-custom) (2018-07-09)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.1.0-custom...v1.1.1-custom)

**Closed issues:**

- Error when updating service LoadBalancerIP [\#23](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/23)
- git clone in instructions shouldn't use git@ [\#20](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/20)
- Missing dependencies on fresh install [\#15](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/15)
- StatefulSet update strategy not propagating [\#13](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/13)
- Update README [\#10](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/10)

**Merged pull requests:**

- Add reclaim policy to storage class [\#29](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/29) ([stigok](https://github.com/stigok))
- Use `Update\(\)` instead of `Patch\(\)` to update svc [\#24](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/24) ([sl1pm4t](https://github.com/sl1pm4t))
- Use HTTPS clone link instead of SSH [\#21](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/21) ([stigok](https://github.com/stigok))
- Add missing dependencies [\#18](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/18) ([stigok](https://github.com/stigok))
- Add instruction on terraform init [\#17](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/17) ([stigok](https://github.com/stigok))
- Fix Ingress update error: [\#16](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/16) ([sl1pm4t](https://github.com/sl1pm4t))
- Fix StatefulSet update strategy [\#14](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/14) ([sl1pm4t](https://github.com/sl1pm4t))

## [v1.1.0-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.1.0-custom) (2018-05-02)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.8-custom...v1.1.0-custom)

**Closed issues:**

- Add field "revision\_history\_limit" to Deployment schema [\#7](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/7)

**Merged pull requests:**

- CronJob resource [\#12](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/12) ([sl1pm4t](https://github.com/sl1pm4t))
- Support Kubernetes v1.9.0 and apps/v1 API [\#11](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/11) ([sl1pm4t](https://github.com/sl1pm4t))
- Add `revision\_history\_limit` to Deployment [\#8](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/8) ([sl1pm4t](https://github.com/sl1pm4t))

## [v1.0.8-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.8-custom) (2018-03-13)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.7-custom...v1.0.8-custom)

**Closed issues:**

- Using imagePullSecrets to pull container from a private repository [\#5](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/5)
- Updating internal annotation fails. [\#2](https://github.com/sl1pm4t/terraform-provider-kubernetes/issues/2)

**Merged pull requests:**

- Set DeletionPolicyForeground on deployment delete [\#4](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/4) ([sl1pm4t](https://github.com/sl1pm4t))

## [v1.0.7-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.7-custom) (2018-02-19)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.6-custom...v1.0.7-custom)

## [v1.0.6-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.6-custom) (2018-02-01)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.5-custom...v1.0.6-custom)

## [v1.0.5-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.5-custom) (2018-01-15)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.4-custom...v1.0.5-custom)

## [v1.0.4-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.4-custom) (2017-12-14)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.3-custom...v1.0.4-custom)

## [v1.0.3-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.3-custom) (2017-10-29)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.2-custom...v1.0.3-custom)

## [v1.0.2-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.2-custom) (2017-10-26)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.1-custom...v1.0.2-custom)

**Merged pull requests:**

- updating deployment structure utility functions to use proper schema … [\#1](https://github.com/sl1pm4t/terraform-provider-kubernetes/pull/1) ([TaysirTayyab](https://github.com/TaysirTayyab))

## [v1.0.1-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.1-custom) (2017-10-26)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v1.0.0...v1.0.1-custom)

## [v1.0.0](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v1.0.0) (2017-08-18)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v0.1.4-custom...v1.0.0)

## [v0.1.4-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v0.1.4-custom) (2017-08-13)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v0.1.3-custom...v0.1.4-custom)

## [v0.1.3-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v0.1.3-custom) (2017-08-11)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v0.1.2-custom...v0.1.3-custom)

## [v0.1.2-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v0.1.2-custom) (2017-08-11)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v0.1.1-custom...v0.1.2-custom)

## [v0.1.1-custom](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v0.1.1-custom) (2017-08-09)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v0.1.2...v0.1.1-custom)

## [v0.1.2](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v0.1.2) (2017-08-04)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v0.1.1...v0.1.2)

## [v0.1.1](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v0.1.1) (2017-07-05)
[Full Changelog](https://github.com/sl1pm4t/terraform-provider-kubernetes/compare/v0.1.0...v0.1.1)

## [v0.1.0](https://github.com/sl1pm4t/terraform-provider-kubernetes/tree/v0.1.0) (2017-06-20)


\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*