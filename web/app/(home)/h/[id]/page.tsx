'use client'
import IRRemoteThingCard from '@/components/thing/ir-remote'
import { useHomeSelector } from '@/lib/contexts/HomeContext'
import { useUser } from '@/lib/hooks'
import { DeviceKindEnum, IDeviceCard } from '@/lib/types'
import { toSWR } from '@/lib/utils'
import DeviceService from '@/services/DeviceService'
import HomeService from '@/services/HomeService'
import useSWR from 'swr'
import HomeGreating from './home-greating'
import HomeSenceList from './home-sence'
import SkeletonDisplay from './skeleton'
const fetchData = async (id: string) => {
  return await HomeService.findById(id)
}

const HomePage = ({ params: { id: homeId } }: { params: { id: string } }) => {
  const x = DeviceService.list(homeId)
  const { data: devices, error, isLoading } = useSWR(`home-devices`, toSWR(x))

  const { profile } = useUser()
  const renderDeviceCard = (d: IDeviceCard) => {
    if (d.kind === DeviceKindEnum.IRRemote) {
      return <IRRemoteThingCard />
    }
    return (
      <div className="col-span-6 md:col-span-4" key={d.id}>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <img
                className="h-10 w-10 rounded-full"
                src="https://via.placeholder.com/400"
                alt="x"
              />
            </div>
          </div>
          <div className="mt-4">
            <div className="text-base font-medium ">{d.name}</div>
            <div>x</div>
          </div>
        </div>
      </div>
    )
  }
  const loading = useHomeSelector((s) => s.isLoading)

  if (loading) return <SkeletonDisplay />

  return (
    <div className="grid grid-cols-12 gap-4">
      <div className="col-span-12">
        <HomeSenceList />
      </div>
      <div className="col-span-12">
        <HomeGreating />
      </div>
      {!error &&
        !isLoading &&
        devices &&
        devices.map((d) => renderDeviceCard(d))}
    </div>
  )
}
export default HomePage
