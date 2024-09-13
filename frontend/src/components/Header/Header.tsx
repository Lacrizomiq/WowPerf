import Link from "next/link";
import Image from "next/image";

export default function Header() {
  return (
    <header className="bg-[#002440] rounded-xl p-5 flex items-center space-x-5 mb-5">
      <nav className="container mx-auto flex justify-between items-center">
        <Link href="/" className="flex items-center">
          <Image src="/wow.png" alt="WoW Logo" width={40} height={40} />
          <span className="ml-2 text-white font-bold text-xl">WoWPerf</span>
        </Link>
        <ul className="flex items-center space-x-6">
          <li>
            <Link
              href="/characters"
              className="text-white hover:text-blue-300 transition-colors"
            >
              Characters
            </Link>
          </li>
          <li>
            <Link
              href="/dungeons"
              className="text-white hover:text-blue-300 transition-colors"
            >
              Dungeons
            </Link>
          </li>
          <li>
            <Link
              href="/raids"
              className="text-white hover:text-blue-300 transition-colors"
            >
              Raids
            </Link>
          </li>
          <li>
            <Link
              href="/signin"
              className="px-4 py-2 text-white bg-gradient-blue rounded-full hover:bg-blue-700 transition-colors"
            >
              Sign In
            </Link>
          </li>
          <li>
            <Link
              href="/signup"
              className="px-4 py-2 text-white bg-gradient-purple rounded-full hover:bg-purple-600 transition-colors"
            >
              Sign Up
            </Link>
          </li>
        </ul>
      </nav>
    </header>
  );
}
