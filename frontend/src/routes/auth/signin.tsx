import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/auth/signin')({
  component: () => <div>Hello /auth/signin!</div>
})