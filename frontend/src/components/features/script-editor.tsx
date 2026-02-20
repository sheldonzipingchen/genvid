'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Loader2, Sparkles, RefreshCw } from 'lucide-react'

const scriptTemplates = [
  {
    id: 'review',
    name: 'Product Review',
    template: `I've been using {product} for a few weeks now, and honestly... it's been a game-changer!

The thing I love most is how {benefit}. It really makes a difference in my daily routine.

If you're someone who {target_audience}, this is definitely worth checking out. I wish I found it sooner!

Link's in bio if you want to grab one for yourself.`,
  },
  {
    id: 'unboxing',
    name: 'Unboxing Experience',
    template: `Okay, so my {product} just arrived and I HAD to share this with you guys!

Look at this packaging... {packaging_reaction}

Let me show you what's inside... {unboxing_content}

First impressions? This quality is insane for the price. Stay tuned for my full review!`,
  },
  {
    id: 'problem-solution',
    name: 'Problem & Solution',
    template: `If you're struggling with {problem}, you NEED to see this.

I used to deal with this ALL the time. It was so frustrating.

But then I found {product}. {solution_description}

Now? Problem solved. And it only took {timeframe}.

Drop a 🔥 if this helped!`,
  },
  {
    id: 'testimonial',
    name: 'Testimonial',
    template: `So I've been using {product} for {timeframe} now...

And I have to be honest - I was skeptical at first.

But {results}? This actually works!

{specific_benefit}

If you're on the fence, just try it. You can thank me later 😉`,
  },
]

interface ScriptEditorProps {
  productName: string
  productDescription: string
  script: string
  videoDuration: number
  onChange: (script: string) => void
  onDurationChange: (duration: number) => void
}

export function ScriptEditor({ productName, productDescription, script, videoDuration, onChange, onDurationChange }: ScriptEditorProps) {
  const [generating, setGenerating] = useState(false)
  const [selectedTemplate, setSelectedTemplate] = useState<string | null>(null)

  const generateScript = async (templateId?: string) => {
    setGenerating(true)
    
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    const template = scriptTemplates.find(t => t.id === templateId) || scriptTemplates[0]
    
    let generatedScript = template.template
      .replace('{product}', productName || 'this product')
      .replace('{benefit}', productDescription?.slice(0, 50) || 'it saves time')
      .replace('{target_audience}', 'wants to level up their life')
      .replace('{problem}', 'this common issue')
      .replace('{solution_description}', 'It completely changed how I approach this.')
      .replace('{timeframe}', 'just 2 weeks')
      .replace('{results}', 'the results speak for themselves')
      .replace('{specific_benefit}', 'The best part? It\'s super easy to use.')
      .replace('{packaging_reaction}', 'so premium!')
      .replace('{unboxing_content}', 'Wow, look at these details!')
    
    onChange(generatedScript)
    setGenerating(false)
  }

  const characterCount = script.length
  const maxCharacters = 1000

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Create Your Script</h2>
        <p className="text-gray-500">Write your own script or let AI generate one for you</p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        {scriptTemplates.map((template) => (
          <Card
            key={template.id}
            className={`cursor-pointer transition-all hover:shadow-md ${
              selectedTemplate === template.id ? 'ring-2 ring-violet-600' : ''
            }`}
            onClick={() => {
              setSelectedTemplate(template.id)
              generateScript(template.id)
            }}
          >
            <CardContent className="p-4 text-center">
              <p className="font-medium text-gray-900 text-sm">{template.name}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="flex gap-3">
        <Button
          variant="outline"
          onClick={() => generateScript(selectedTemplate || 'review')}
          disabled={generating || !productName}
        >
          {generating ? (
            <Loader2 className="h-4 w-4 animate-spin mr-2" />
          ) : (
            <Sparkles className="h-4 w-4 mr-2" />
          )}
          {generating ? 'Generating...' : 'Generate with AI'}
        </Button>
        <Button
          variant="ghost"
          onClick={() => generateScript(selectedTemplate || 'review')}
          disabled={generating}
        >
          <RefreshCw className="h-4 w-4 mr-2" />
          Regenerate
        </Button>
      </div>

      <div>
        <div className="flex justify-between items-center mb-2">
          <label className="block text-sm font-medium text-gray-700">
            Your Script
          </label>
          <span className={`text-xs ${characterCount > maxCharacters ? 'text-red-500' : 'text-gray-500'}`}>
            {characterCount} / {maxCharacters} characters
          </span>
        </div>
        <textarea
          value={script}
          onChange={(e) => onChange(e.target.value)}
          placeholder="Write your script here or click 'Generate with AI' to create one automatically..."
          rows={10}
          className={`w-full rounded-lg border p-4 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500 ${
            characterCount > maxCharacters ? 'border-red-300' : 'border-gray-200'
          }`}
        />
      </div>

      <div className="bg-violet-50 rounded-lg p-4">
        <h4 className="font-medium text-violet-900 mb-2">💡 Tips for Better Videos</h4>
        <ul className="text-sm text-violet-700 space-y-1">
          <li>• Keep it conversational and natural</li>
          <li>• Focus on benefits, not just features</li>
          <li>• Include a clear call-to-action</li>
          <li>• Aim for 30-60 seconds (75-150 words)</li>
        </ul>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Video Duration
        </label>
        <div className="flex gap-3">
          {[
            { value: 5, label: '5s', desc: 'Quick teaser' },
            { value: 10, label: '10s', desc: 'Standard' },
            { value: 30, label: '30s', desc: 'Full ad' },
          ].map((duration) => (
            <Card
              key={duration.value}
              className={`flex-1 cursor-pointer transition-all ${
                videoDuration === duration.value ? 'ring-2 ring-violet-600' : ''
              }`}
              onClick={() => onDurationChange(duration.value)}
            >
              <CardContent className="p-4 text-center">
                <p className="font-medium text-gray-900">{duration.label}</p>
                <p className="text-xs text-gray-500">{duration.desc}</p>
              </CardContent>
            </Card>
          ))}
        </div>
        {videoDuration === 30 && (
          <p className="text-xs text-amber-600 mt-2">
            ⚠️ 30-second videos will be generated as 3 segments and merged, taking longer to process.
          </p>
        )}
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Video Format
        </label>
        <div className="flex gap-3">
          {[
            { value: '9:16', label: '9:16', desc: 'TikTok / Reels' },
            { value: '1:1', label: '1:1', desc: 'Square' },
          ].map((format) => (
            <Card
              key={format.value}
              className={`flex-1 cursor-pointer transition-all ${
                format.value === '9:16' ? 'ring-2 ring-violet-600' : ''
              }`}
            >
              <CardContent className="p-4 text-center">
                <div className={`mx-auto mb-2 border-2 rounded ${
                  format.value === '9:16' 
                    ? 'w-6 h-10 border-violet-600' 
                    : 'w-8 h-8 border-gray-300'
                }`} />
                <p className="font-medium text-gray-900">{format.label}</p>
                <p className="text-xs text-gray-500">{format.desc}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  )
}
