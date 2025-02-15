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
