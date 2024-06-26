"use client"

import { usePathname } from "next/navigation"
import { TfiMoney } from "react-icons/tfi"
import { IoHome } from "react-icons/io5"
import Image from "next/image"
import Link from "next/link"

export function Sidebar() {
  const path = usePathname()
  return (
    <div className="bg-dashboard sticky top-0 h-screen border-r border-r-white border-opacity-30 text-white">
      <div className="flex flex-col">
        <Link
          className="mb-5 mr-10 flex flex-row items-center gap-x-2 p-5 text-3xl font-bold"
          href="/"
        >
          <Image src="/logos/box.svg" alt="logo-box" width={40} height={40} />
          BlockFeed
        </Link>
        <div className="flex flex-col gap-y-5">
          {items.map((item, i) => (
            <Link
              className={"ml-5 mr-20".concat(
                path === item.route
                  ? ""
                  : " ".concat(
                      "opacity-50 transition-all ease-linear hover:opacity-100",
                    ),
              )}
              key={i}
              href={item.route}
            >
              <div className="flex flex-row items-center gap-x-4">
                {item.icon}
                {item.name}
              </div>
            </Link>
          ))}
        </div>
      </div>
    </div>
  )
}

const items = [
  {
    route: "/",
    name: "Dashboard",
    icon: <IoHome />,
  },
  {
    route: "/billing",
    name: "Billing",
    icon: <TfiMoney />,
  },
]
