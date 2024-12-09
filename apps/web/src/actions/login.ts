"use server"

export async function login(data: { email: string; password: string }) {
  try {
    const response = await fetch(
      process.env.NEXT_PUBLIC_BACKEND_URL + "/api/v1/auth/login",
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    )

    if (!response.ok) {
      const { error } = await response.json()
      throw new Error(error || "Login failed")
    }

    const tokens = await response.json()

    console.log(">>>> TOKEN", tokens.access_token)
    console.log(">>>> TOKEN", tokens.refresh_token)

    return tokens
  } catch (error) {
    console.log(error)
    return null
  }
}
