# https://trivy.dev/v0.47/docs/configuration/filtering/#by-open-policy-agent

package trivy

import data.lib.trivy

default ignore = false

# ----------------------
# Distroless OS packages
# ----------------------

ignore {
    input.PkgName == "base-files"
    input.Name == "GPL-2.0-or-later"
}

ignore {
    input.PkgName == "base-files"
    input.Name == "verbatim"
}

ignore {
    input.PkgName == "ca-certificates"
    {"GPL-2.0-or-later", "GPL-2.0-only", "MPL-2.0"}[input.Name]
}

ignore {
    input.PkgName == "netbase"
    input.Name == "GPL-2.0-only"
}

ignore {
    input.PkgName == "tzdata"
    input.Name == "public-domain"
}

ignore {
    input.PkgName == "tzdata-legacy"
    input.Name == "public-domain"
}

ignore {
    input.PkgName = "media-types"
    input.Name = "ad-hoc"
}

# ---------------------------------
# Copyleft licensed Go dependencies
# ---------------------------------

# MPL-2.0 - changes must retain license, cannot mix licenses in same file
ignore {
    input.PkgName = "github.com/go-sql-driver/mysql"
    input.Name = "MPL-2.0"
}

ignore {
    input.PkgName = "github.com/hashicorp/go-immutable-radix"
    input.Name = "MPL-2.0"
}

ignore {
    input.PkgName = "github.com/hashicorp/go-memdb"
    input.Name = "MPL-2.0"
}

ignore {
    input.PkgName = "github.com/hashicorp/golang-lru"
    input.Name = "MPL-2.0"
}
