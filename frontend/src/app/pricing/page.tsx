'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { MarketingLayout } from '@/components/layout/marketing-layout'
import { useAuthStore } from '@/stores/auth'
import { useStripe } from '@/hooks/use-stripe'
import { Check, Zap, Building, Crown, Loader2 } from 'lucide-react'

const plans = [
  {
    id: 'free',
    name: 'Free',
    price: '$0',
    period: 'forever',
    description: 'Try out Genvid with limited features',
    icon: Zap,
    features: [
      '3 videos per month',
      'Basic avatars',
      '720p export',
      'Watermark included',
      'Community support',
    ],
    cta: 'Get Started',
    ctaLink: '/register',
    popular: false,
    requiresPayment: false,
  },
  {
    id: 'starter_monthly',
    name: 'Starter',
    price: '$19',
    period: '/month',
    description: 'Perfect for small e-commerce sellers',
    icon: Zap,
    features: [
      '15 videos per month',
      'All basic avatars',
      '1080p export',
      'No watermark',
      'Email support',
      'Script templates',
    ],
    cta: 'Start Free Trial',
    ctaLink: '/register?plan=starter',
    popular: true,
    requiresPayment: true,
  },
  {
    id: 'pro_monthly',
    name: 'Pro',
    price: '$49',
    period: '/month',
    description: 'For growing businesses and agencies',
    icon: Building,
    features: [
      '50 videos per month',
      'All avatars including premium',
      '1080p export',
      'No watermark',
      'Priority support',
      'Custom scripts',
      'Multi-language',
      'Brand kit',
    ],
    cta: 'Start Free Trial',
    ctaLink: '/register?plan=pro',
    popular: false,
    requiresPayment: true,
  },
  {
    id: 'business_monthly',
    name: 'Business',
    price: '$99',
    period: '/month',
    description: 'For agencies and high-volume sellers',
    icon: Crown,
    features: [
      '150 videos per month',
      'All avatars + custom avatars',
      '4K export',
      'No watermark',
      'Dedicated support',
      'API access',
      'Team collaboration',
      'White-label option',
      'Analytics dashboard',
    ],
    cta: 'Contact Sales',
    ctaLink: '/contact?plan=business',
    popular: false,
    requiresPayment: true,
  },
]

const faqs = [
  {
    question: 'How long does it take to generate a video?',
    answer: 'Most videos are generated within 2-5 minutes. Complex videos with longer scripts may take up to 10 minutes.',
  },
  {
    question: 'Can I use the videos commercially?',
    answer: 'Yes! All videos you generate can be used for your business, ads, and social media without any restrictions.',
  },
  {
    question: 'What happens if I run out of credits?',
    answer: 'You can upgrade your plan anytime or purchase additional credits. Unused credits roll over to the next month.',
  },
  {
    question: 'Can I cancel my subscription?',
    answer: 'Yes, you can cancel anytime. Your subscription remains active until the end of the billing period.',
  },
  {
    question: 'Do you offer refunds?',
    answer: 'We offer a 7-day money-back guarantee. If you\'re not satisfied, contact support for a full refund.',
  },
]

export default function PricingPage() {
  const router = useRouter()
  const { isAuthenticated } = useAuthStore()
  const { createCheckoutSession, loading, error } = useStripe()
  const [processingPlan, setProcessingPlan] = useState<string | null>(null)

  const handlePlanSelect = async (plan: typeof plans[0]) => {
    if (!plan.requiresPayment) {
      router.push(plan.ctaLink)
      return
    }

    if (!isAuthenticated) {
      router.push(`/register?plan=${plan.id}`)
      return
    }

    setProcessingPlan(plan.id)
    
    const result = await createCheckoutSession({
      planId: plan.id,
    })

    if (!result) {
      setProcessingPlan(null)
    }
  }

  return (
    <MarketingLayout>
      <section className="py-24 bg-white">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h1 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4">
              Simple, Transparent Pricing
            </h1>
            <p className="text-lg text-gray-600">
              Start free and scale as you grow. No hidden fees, no surprises.
            </p>
            {error && (
              <p className="mt-4 text-red-600 text-sm">{error}</p>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {plans.map((plan) => {
              const Icon = plan.icon
              const isLoading = processingPlan === plan.id
              
              return (
                <Card
                  key={plan.name}
                  className={`relative ${plan.popular ? 'ring-2 ring-violet-600 shadow-xl' : ''}`}
                >
                  {plan.popular && (
                    <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                      <span className="bg-violet-600 text-white text-xs font-medium px-3 py-1 rounded-full">
                        Most Popular
                      </span>
                    </div>
                  )}
                  <CardHeader className="text-center pb-4">
                    <div className="h-12 w-12 mx-auto mb-4 rounded-full bg-violet-100 flex items-center justify-center">
                      <Icon className="h-6 w-6 text-violet-600" />
                    </div>
                    <CardTitle className="text-xl">{plan.name}</CardTitle>
                    <div className="mt-2">
                      <span className="text-4xl font-bold text-gray-900">{plan.price}</span>
                      <span className="text-gray-500">{plan.period}</span>
                    </div>
                    <CardDescription className="mt-2">{plan.description}</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <ul className="space-y-2">
                      {plan.features.map((feature) => (
                        <li key={feature} className="flex items-start gap-2 text-sm text-gray-600">
                          <Check className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                          {feature}
                        </li>
                      ))}
                    </ul>
                    <Button
                      className="w-full"
                      variant={plan.popular ? 'default' : 'outline'}
                      onClick={() => handlePlanSelect(plan)}
                      disabled={loading && isLoading}
                    >
                      {isLoading ? (
                        <>
                          <Loader2 className="h-4 w-4 animate-spin mr-2" />
                          Processing...
                        </>
                      ) : (
                        plan.cta
                      )}
                    </Button>
                  </CardContent>
                </Card>
              )
            })}
          </div>

          <div className="mt-16 text-center">
            <p className="text-gray-500 mb-4">Need more? Enterprise plans available for high-volume needs.</p>
            <Button variant="outline" asChild>
              <Link href="/contact">Contact Sales</Link>
            </Button>
          </div>
        </div>
      </section>

      <section className="py-24 bg-gray-50">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-gray-900 text-center mb-12">
            Frequently Asked Questions
          </h2>
          <div className="space-y-6">
            {faqs.map((faq, index) => (
              <Card key={index}>
                <CardContent className="p-6">
                  <h3 className="font-semibold text-gray-900 mb-2">{faq.question}</h3>
                  <p className="text-gray-600">{faq.answer}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      <section className="py-24 bg-gradient-to-r from-violet-600 to-purple-600">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            Start Creating Videos Today
          </h2>
          <p className="text-lg text-violet-100 mb-8 max-w-2xl mx-auto">
            Join thousands of e-commerce sellers using Genvid to create viral UGC content.
          </p>
          <Button size="xl" variant="secondary" asChild>
            <Link href="/register">Get Started Free</Link>
          </Button>
        </div>
      </section>
    </MarketingLayout>
  )
}
