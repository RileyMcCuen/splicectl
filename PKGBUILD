# The splicectl cli is used to manage features of a SpliceDB Cluster running on Kubernetes.

# Maintainer: Your Name <blo@splicemachine.com>
pkgname=("splicectl")
pkgver=v0.1.1
pkgrel=1
epoch=
pkgdesc="cli is used to manage features of a SpliceDB Cluster running on Kubernetes."
arch=('x86_64')
url="https://github.com/splicemachine/splicectl"
license=('GPL3')
groups=()
optdepends=()
provides=()
conflicts=()
replaces=()
backup=()
options=()
changelog=
source=("$pkgname-$pkgver::git+https://github.com/splicemachine/splicectl/")
noextract=()
md5sums=('SKIP')
validpgpkeys=()

package() {
	cd "$pkgname-$pkgver"
	install -Dm755 "$srcdir/$pkgname-$pkgver/splicectl" "$pkgdir/usr/bin/splicectl"
}
