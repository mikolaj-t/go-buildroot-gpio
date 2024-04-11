################################################################################
#
# gogpio
#
################################################################################
GOGPIO_VERSION = 0.10
GOGPIO_SITE = $(call github,mikolaj-t,go-buildroot-gpio,$(GOGPIO_VERSION))

$(eval $(golang-package))