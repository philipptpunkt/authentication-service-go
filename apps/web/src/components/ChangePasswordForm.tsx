"use client"

import { SubmitHandler, useForm } from "react-hook-form"

type Inputs = {
  currentPassword: string
  newPassword: string
}

interface LoginFormProps {
  onSubmit: SubmitHandler<Inputs>
}

export function ChangePasswordForm({ onSubmit }: LoginFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Inputs>()

  return (
    <form
      onSubmit={handleSubmit((data) =>
        onSubmit({
          currentPassword: data.currentPassword,
          newPassword: data.newPassword,
        })
      )}
    >
      <input
        className="mt-4 px-4 py-2 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
        type="password"
        placeholder="Current password"
        {...register("currentPassword", { required: true })}
      />

      <input
        className="mt-4 px-4 py-2 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
        type="password"
        placeholder="New password"
        {...register("newPassword", { required: true })}
      />
      {errors.currentPassword && <span>Old Password is required</span>}
      {errors.newPassword && <span>New password is required</span>}

      <input
        type="submit"
        className="bg-black text-white py-2 w-full mt-4 rounded-md cursor-pointer text-lg"
        value="Set new password"
      />
    </form>
  )
}
