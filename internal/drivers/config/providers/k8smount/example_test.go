package k8smount_test

import (
	"log"

	"github.com/knadh/koanf/v2"

	"github.com/mattdowdell/sandbox/internal/drivers/config/providers/k8smount"
)

func ExampleK8SMount_Read() {
	k := koanf.New(".")
	provider := k8smount.Provider("/path/to/mount/", "." /*delimiter*/)

	if err := k.Load(provider, nil /*parser*/); err != nil {
		log.Fatal("config load failed:", err)
		return
	}

	log.Println("config loaded: example=", k.Get("example"))
}

func ExampleK8SMount_Watch() {
	k := koanf.New(".")
	provider := k8smount.Provider("/path/to/mount/", "." /*delimiter*/)

	if err := k.Load(provider, nil /*parser*/); err != nil {
		log.Fatal("config load failed:", err)
		return
	}

	log.Println("config loaded: example=", k.Get("example"))

	if err := provider.Watch(func(err error) {
		if err != nil {
			log.Println("watch stopping:", err)
			return
		}

		if err := k.Load(provider, nil /*parser*/); err != nil {
			log.Println("config load failed:", err)

			if err := provider.Unwatch(); err != nil {
				log.Println("provider unwatch error:", err)
			}

			return
		}

		log.Println("config loaded: example=", k.Get("example"))
	}); err != nil {
		log.Fatal("watch failed:", err)
	}

	// wait forever
	select {}
}
