################################################################################
#
# gogpio
#
################################################################################
GOGPIO_VERSION = 0.8
GOGPIO_SITE = $(call github,mikolaj-t,go-buildroot-gpio,$(GOGPIO_VERSION))

$(eval $(golang-package))