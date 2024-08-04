import Link from "next/link";
import Image from "next/image";

export default function Header() {
  return (
    <header className="bg-gray-800 py-4">
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
              className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 transition duration-300 border-2 border-blue-300"
            >
              Sign Up
            </Link>
          </li>
        </ul>
      </nav>
    </header>
  );
}
