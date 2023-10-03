pkgname=synctropy
pkgver=0.1.0
pkgrel=1
pkgdesc="A wrapper for management and syncing of crates via syncing utilities like unison and rsync using hooks, with template support."
arch=('any')
url="https://github.com/fearlessdots/synctropy"
license=('GPL3')
depends=('glibc' 'gcc-libs')
makedepends=('go')
source=("${pkgname}-${pkgver}.tar.gz::${url}/archive/refs/tags/v${pkgver}.tar.gz")
sha256sums=('SKIP')

build() {
	cd "$pkgname-$pkgver"
	make build
}

package() {
	# Create directories for shell autocompletion
	echo "Creating directories for shell autocompletion"
	mkdir -p ${pkgdir}/usr/share/bash-completion/completions ${pkgdir}/usr/share/zsh/site-functions \
		${pkgdir}/usr/share/fish/vendor_completions.d

	cd "$pkgname-$pkgver"
	make DESTDIR=$pkgdir install
}
