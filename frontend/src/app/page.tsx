import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { MarketingLayout } from '@/components/layout/marketing-layout'
import { 
  Video, 
  Sparkles, 
  Zap, 
  Globe, 
  Check, 
  Play,
  ArrowRight,
  Star
} from 'lucide-react'

const features = [
  {
    icon: Zap,
    title: 'Lightning Fast',
    description: 'Generate professional UGC videos in just 2-5 minutes. No waiting weeks for creators.',
  },
  {
    icon: Sparkles,
    title: 'AI-Powered',
    description: 'Our AI avatars are trained on real creators for authentic, engaging content.',
  },
  {
    icon: Globe,
    title: 'Multi-Language',
    description: 'Reach global audiences with support for 10+ languages and localized avatars.',
  },
  {
    icon: Video,
    title: 'All Formats',
    description: 'Export in 9:16, 1:1, or 16:9 ratios. Perfect for TikTok, Reels, and Shorts.',
  },
]

const steps = [
  {
    number: '01',
    title: 'Add Your Product',
    description: 'Enter your product URL or upload images. We automatically extract product details.',
  },
  {
    number: '02',
    title: 'Choose Avatar & Script',
    description: 'Select from 100+ AI avatars and let our AI generate the perfect script for your product.',
  },
  {
    number: '03',
    title: 'Generate & Download',
    description: 'Click generate and get your UGC video in minutes. Download and use immediately.',
  },
]

const testimonials = [
  {
    content: "Genvid cut our video production costs by 80%. We went from 5 videos/month to 50+.",
    author: "Sarah Chen",
    role: "E-commerce Manager",
    company: "TechStyle",
  },
  {
    content: "The AI avatars look so authentic. Our TikTok engagement increased by 3x in the first month.",
    author: "Michael Torres",
    role: "Marketing Director",
    company: "FitGear Co.",
  },
  {
    content: "Finally, a tool that understands e-commerce. The multi-language feature opened up EU markets for us.",
    author: "Emma Watson",
    role: "Founder",
    company: "GlowBeauty",
  },
]

export default function HomePage() {
  return (
    <MarketingLayout>
      <section className="relative overflow-hidden bg-gradient-to-b from-violet-50 via-white to-white">
        <div className="absolute inset-0 bg-grid-slate-100 [mask-image:linear-gradient(0deg,white,rgba(255,255,255,0))]" />
        
        <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-24 sm:py-32">
          <div className="text-center max-w-4xl mx-auto">
            <div className="inline-flex items-center gap-2 rounded-full bg-violet-100 px-4 py-1.5 text-sm font-medium text-violet-700 mb-8">
              <Sparkles className="h-4 w-4" />
              <span>AI-Powered UGC Video Generator</span>
            </div>
            
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight text-gray-900 mb-6">
              Turn Products into{' '}
              <span className="bg-gradient-to-r from-violet-600 to-purple-600 bg-clip-text text-transparent">
                Viral Videos
              </span>{' '}
              in Minutes
            </h1>
            
            <p className="text-lg sm:text-xl text-gray-600 mb-10 max-w-2xl mx-auto">
              Create authentic UGC videos for TikTok, Reels, and Shorts. 
              10x faster, 80% cheaper than traditional creators.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Button size="xl" asChild>
                <Link href="/register">
                  Start Creating Free
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button size="xl" variant="outline" asChild>
                <Link href="/demo" className="flex items-center gap-2">
                  <Play className="h-4 w-4" />
                  Watch Demo
                </Link>
              </Button>
            </div>
            
            <div className="mt-10 flex items-center justify-center gap-8 text-sm text-gray-500">
              <div className="flex items-center gap-2">
                <Check className="h-4 w-4 text-green-500" />
                <span>No credit card required</span>
              </div>
              <div className="flex items-center gap-2">
                <Check className="h-4 w-4 text-green-500" />
                <span>3 free videos</span>
              </div>
              <div className="flex items-center gap-2">
                <Check className="h-4 w-4 text-green-500" />
                <span>Cancel anytime</span>
              </div>
            </div>
          </div>
          
          <div className="mt-16 relative max-w-5xl mx-auto">
            <div className="absolute -inset-4 bg-gradient-to-r from-violet-500 to-purple-500 rounded-2xl blur-2xl opacity-20" />
            <div className="relative bg-gray-900 rounded-2xl shadow-2xl overflow-hidden">
              <div className="flex items-center gap-2 px-4 py-3 bg-gray-800 border-b border-gray-700">
                <div className="w-3 h-3 rounded-full bg-red-500" />
                <div className="w-3 h-3 rounded-full bg-yellow-500" />
                <div className="w-3 h-3 rounded-full bg-green-500" />
              </div>
              <div className="aspect-video bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center">
                <div className="text-center p-8">
                  <div className="w-24 h-24 mx-auto mb-6 rounded-full bg-gradient-to-br from-violet-500 to-purple-500 flex items-center justify-center">
                    <Video className="h-12 w-12 text-white" />
                  </div>
                  <p className="text-gray-400 text-lg">Video Preview</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="py-12 bg-white border-y border-gray-100">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-8">
            <div className="flex items-center gap-4">
              <div className="flex -space-x-2">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div
                    key={i}
                    className="h-10 w-10 rounded-full bg-gradient-to-br from-violet-400 to-purple-400 border-2 border-white"
                  />
                ))}
              </div>
              <div>
                <p className="text-sm font-medium text-gray-900">Trusted by 5,000+ sellers</p>
                <div className="flex items-center gap-1">
                  {[1, 2, 3, 4, 5].map((i) => (
                    <Star key={i} className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                  ))}
                  <span className="ml-2 text-sm text-gray-500">4.9/5 rating</span>
                </div>
              </div>
            </div>
            <div className="flex items-center gap-8 text-gray-400">
              {['TikTok Shop', 'Shopify', 'Amazon', 'WooCommerce'].map((brand) => (
                <span key={brand} className="text-sm font-medium">{brand}</span>
              ))}
            </div>
          </div>
        </div>
      </section>

      <section className="py-24 bg-white">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
              Everything You Need for UGC Videos
            </h2>
            <p className="text-lg text-gray-600">
              From product to viral video in just 3 clicks. No camera, no crew, no editing skills required.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature) => (
              <Card key={feature.title} className="border-0 shadow-lg hover:shadow-xl transition-shadow">
                <CardContent className="p-6">
                  <div className="h-12 w-12 rounded-lg bg-violet-100 flex items-center justify-center mb-4">
                    <feature.icon className="h-6 w-6 text-violet-600" />
                  </div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">{feature.title}</h3>
                  <p className="text-gray-600">{feature.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      <section className="py-24 bg-gray-50">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
              Create Videos in 3 Simple Steps
            </h2>
            <p className="text-lg text-gray-600">
              No learning curve. No complicated software. Just results.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {steps.map((step) => (
              <div key={step.number} className="relative">
                <div className="text-6xl font-bold text-violet-100 mb-4">{step.number}</div>
                <h3 className="text-xl font-semibold text-gray-900 mb-2">{step.title}</h3>
                <p className="text-gray-600">{step.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="py-24 bg-white">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
              Loved by E-commerce Sellers
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {testimonials.map((testimonial) => (
              <Card key={testimonial.author} className="border-0 shadow-lg">
                <CardContent className="p-6">
                  <div className="flex items-center gap-1 mb-4">
                    {[1, 2, 3, 4, 5].map((i) => (
                      <Star key={i} className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                    ))}
                  </div>
                  <p className="text-gray-600 mb-6">&ldquo;{testimonial.content}&rdquo;</p>
                  <div>
                    <p className="font-semibold text-gray-900">{testimonial.author}</p>
                    <p className="text-sm text-gray-500">{testimonial.role}, {testimonial.company}</p>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      <section className="py-24 bg-gradient-to-r from-violet-600 to-purple-600">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            Ready to Create Your First Video?
          </h2>
          <p className="text-lg text-violet-100 mb-8 max-w-2xl mx-auto">
            Join thousands of e-commerce sellers who are already creating viral UGC content with Genvid.
          </p>
          <Button size="xl" variant="secondary" asChild>
            <Link href="/register">
              Get Started Free
              <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
          </Button>
        </div>
      </section>
    </MarketingLayout>
  )
}
