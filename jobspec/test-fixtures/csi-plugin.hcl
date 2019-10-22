job "binstore-storagelocker" {
  group "binsl" {
    task "binstore" {
      driver = "docker"

      csi_plugin {
<<<<<<< HEAD
        id        = "org.hashicorp.csi"
        type      = "monolith"
        mount_dir = "/csi/test"
=======
        plugin_id        = "org.hashicorp.csi"
        plugin_type      = "monolith"
        plugin_mount_dir = "/csi/test"
>>>>>>> cfg: Allow specifying the `csi_plugin` stanza
      }
    }
  }
}
