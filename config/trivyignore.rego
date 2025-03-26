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
    input.PkgName == "netbase"
    input.Name == "GPL-2.0-only"
}

ignore {
    input.PkgName == "tzdata"
    input.Name == "public-domain"
}

# ---------------------------------
# Copyleft licensed Go dependencies
# ---------------------------------

# MPL-2.0 - changes must retain license, cannot mix licenses in same file
ignore {
    input.PkgName = "github.com/go-sql-driver/mysql"
    input.Name = "MPL-2.0"
}
