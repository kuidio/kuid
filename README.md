![KUID logo](https://kuidio.github.io/docs/assets/logos/KUID-logo-100x123.png)

# Kuid (Kubernetes Identities)

[![Discord](https://img.shields.io/discord/1234818321833136199?style=flat-square&label=discord&logo=discord&color=00c9ff&labelColor=bec8d2)](https://discord.gg/hXt4sfUs6V)


## What is this?

`Kuid` is a cloud-native application that extends the Kubernetes API, dedicated to managing resources (inventory, IP, VLAN, AS, etc.) within your Kubernetes environments. By leveraging Kubernetes-native architecture and customizable fields, Kuid facilitates streamlined resource organization and tracking, offering notable features such as robust IP Address Management (IPAM) capabilities for efficient allocation and oversight of IP resources. Additionally, Kuid provides sophisticated infrastructure management functionalities, empowering users to organize and manage various infrastructure components within a structured hierarchy.

One of `Kuid's` standout features is its flexibility in resource management, allowing users to define resources statically or discover them dynamically. This dynamic discovery capability enables seamless integration with applications built on top of the Kuid API, empowering users to automate resource provisioning and management tasks effectively.

Moreover, `Kuid` introduces the concept of dynamic resource claiming based on selectors, enabling efficient resource allocation based on specific criteria. With these capabilities, Kuid provides a robust and flexible framework for managing resources effectively, whether in traditional or cloud-native environments. By offering comprehensive inventory management and precise resource identification, Kuid empowers users to optimize infrastructure operations and streamline resource provisioning workflows.

Leveraging Kubernetes-native architecture and customizable fields, Kuid serves as a cornerstone for automation. Seamlessly integrating into automation workflows, `Kuid` emerges as a key building block for driving efficiency and scalability in automating your infrastructure.

For more information consult to [https://kuidio.github.io/docs/](https://kuidio.github.io/docs/)

## How to engage?

* join [our discord server](https://discord.gg/fH35bmcTU9)
* raise [an issue](https://github.com/kuidio/kuid/issues)
* create a [pull request](https://github.com/kuidio/kuid/pulls)
* [how to contribute](CONTRIBUTING.md)

## License and governance

Code in the KUID repositories licensed with [Apache License 2.0](LICENSE.md). At the moment the project is governed by the benevolent dictatorship of @henderiw @steiler @karimra and @hansthienpondt . On the long run we plan to move to a meritocracy based governance model.

## Badges

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/kuidio/kuid/badge)](https://scorecard.dev/viewer/?uri=github.com/kuidio/kuid)



/Users/henderiw/go/bin/go-to-protobuf --go-header-file hack/boilerplate.go.txt --packages ./apis/common/v1alpha1 --apimachinery-packages -k8s.io/apimachinery/pkg/api/resource,-k8s.io/apimachinery/pkg/runtime/schema,-k8s.io/apimachinery/pkg/runtime,-k8s.io/apimachinery/pkg/apis/meta/v1


/Users/henderiw/go/bin/go-to-protobuf --go-header-file hack/boilerplate.go.txt --packages ./apis/config/v1alpha1 --apimachinery-packages -k8s.io/apimachinery/pkg/api/resource,-k8s.io/apimachinery/pkg/runtime/schema,-k8s.io/apimachinery/pkg/runtime,-k8s.io/apimachinery/pkg/apis/meta/v1