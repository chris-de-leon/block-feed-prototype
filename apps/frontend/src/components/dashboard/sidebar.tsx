import { IoHome, IoSettingsSharp } from "react-icons/io5"
import { IoStatsChart } from "react-icons/io5"
import { TfiMoney } from "react-icons/tfi"
import { useRouter } from "next/router"
import Image from "next/image"
import Link from "next/link"

export function Sidebar() {
  const router = useRouter()
  return (
    <div className="sticky top-0 h-screen border-r border-r-white border-opacity-30 bg-dashboard text-white">
      <div className="flex flex-col">
        <Link
          className="mb-5 mr-10 flex flex-row items-center gap-x-2 p-5 text-3xl font-bold"
          href="/dashboard/"
        >
          <Image src="/logos/box.svg" alt="logo-box" width={40} height={40} />
          BlockFeed
        </Link>
        <div className="flex flex-col gap-y-5">
          {items.map((item, i) => {
            return (
              <Link
                className={"ml-5 mr-20".concat(
                  router.route === item.route
                    ? ""
                    : " opacity-50 transition-all ease-linear hover:opacity-100",
                )}
                key={i}
                href={item.route}
              >
                <div className="flex flex-row items-center gap-x-4">
                  {item.icon}
                  {item.name}
                </div>
              </Link>
            )
          })}
        </div>
      </div>
    </div>
  )
}

const items = [
  {
    route: "/dashboard",
    name: "Dashboard",
    icon: <IoHome />,
  },
  {
    route: "/dashboard/monitoring",
    name: "Monitoring",
    icon: <IoStatsChart />,
  },
  {
    route: "/dashboard/billing",
    name: "Billing",
    icon: <TfiMoney />,
  },
  {
    route: "/dashboard/settings",
    name: "Settings",
    icon: <IoSettingsSharp />,
  },
]
