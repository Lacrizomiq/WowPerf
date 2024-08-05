import Link from "next/link";
import Image from "next/image";

export default function Header() {
  return (
    <header className="bg-gradient-dark py-4 px-4">
      <nav className="container mx-auto flex justify-between items-center">
        <Link href="/">
          <Image src="/wow.png" alt="WoW Logo" width={50} height={50} />
        </Link>
        <ul className="flex space-x-6">
          <li>
            <Link href="/characters" className="hover:text-blue-400">
              Characters
            </Link>
          </li>
          <li>
            <Link href="/dungeons" className="hover:text-blue-400">
              Dungeons
            </Link>
          </li>
          <li>
            <Link href="/raids" className="hover:text-blue-400">
              Raids
            </Link>
          </li>
          <li>
            <Link href="/signin" className="text-white hover:text-blue-400">
              Sign In
            </Link>
          </li>
          <li>
            <Link
              href="/signup"
              className="w-full md:w-auto bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition duration-300 glow-effect"
            >
              Sign Up
            </Link>
          </li>
        </ul>
      </nav>
    </header>
  );
}
