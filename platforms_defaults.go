package main

func defaultDarwin() []Platform {
	return []Platform{
		{OS_DARWIN, ARCH_386},
		{OS_DARWIN, ARCH_AMD64},
		//{OS_DARWIN, ARCH_ARM},
		//{OS_DARWIN, ARCH_ARM64},
	}
}

func defaultDragonfly() []Platform {
	return []Platform{
		{OS_DRAGONFLY, ARCH_AMD64},
	}
}

func defaultFreeBSD() []Platform {
	return []Platform{
		{OS_FREEBSD, ARCH_386},
		{OS_FREEBSD, ARCH_AMD64},
		//{OS_FREEBSD, ARCH_ARM},
	}
}

func defaultLinux() []Platform {
	return []Platform{
		{OS_LINUX, ARCH_386},
		{OS_LINUX, ARCH_AMD64},
		//{OS_LINUX, ARCH_ARM},
		//{OS_LINUX, ARCH_ARM64},
		{OS_LINUX, ARCH_PPC64},
		{OS_LINUX, ARCH_PPC64LE},
		{OS_LINUX, ARCH_MIPS},
		{OS_LINUX, ARCH_MIPSLE},
		{OS_LINUX, ARCH_MIPS64},
		{OS_LINUX, ARCH_MIPS64LE},
	}
}

func defaultNetBSD() []Platform {
	return []Platform{
		{OS_NETBSD, ARCH_386},
		{OS_NETBSD, ARCH_AMD64},
		//{OS_NETBSD, ARCH_ARM},
	}
}

func defaultOpenBSD() []Platform {
	return []Platform{
		{OS_OPENBSD, ARCH_386},
		{OS_OPENBSD, ARCH_AMD64},
		//{OS_OPENBSD, ARCH_ARM},
	}
}

func defaultPlan9() []Platform {
	return []Platform{
		{OS_PLAN9, ARCH_386},
		{OS_PLAN9, ARCH_AMD64},
	}
}

func defaultSolaris() []Platform {
	return []Platform{
		{OS_SOLARIS, ARCH_AMD64},
	}
}

func defaultWindows() []Platform {
	return []Platform{
		{OS_WINDOWS, ARCH_386},
		{OS_WINDOWS, ARCH_AMD64},
	}
}
