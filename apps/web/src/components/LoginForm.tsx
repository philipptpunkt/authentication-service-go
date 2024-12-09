"use client"

import { SubmitHandler, useForm } from "react-hook-form"

type Inputs = {
  email: string
  password: string
}

interface LoginFormProps {
  onSubmit: SubmitHandler<Inputs>
}

export function LoginForm({ onSubmit }: LoginFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Inputs>()

  return (
    <form
      onSubmit={handleSubmit((data) =>
        onSubmit({ email: data.email, password: data.password })
      )}
    >
      <input
        className="mt-4 px-4 py-2 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
        type="email"
        placeholder="email"
        {...register("email", { required: true })}
      />

      <input
        className="mt-4 px-4 py-2 block w-full rounded-md border border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
        type="password"
        placeholder="password"
        {...register("password", { required: true })}
      />
      {errors.email && <span>Email is required</span>}
      {errors.password && <span>Password is required</span>}

      <input
        type="submit"
        className="bg-black text-white py-2 w-full mt-4 rounded-md cursor-pointer text-lg"
        value="Log in"
      />
    </form>
  )
}
