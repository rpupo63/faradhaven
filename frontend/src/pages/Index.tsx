import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import { LogIn, Sparkles, BookOpen, Users, Swords, ArrowRight } from 'lucide-react';

export default function Index() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  return (
    <div className="w-full space-y-16 py-12 px-4 md:px-8 max-w-5xl mx-auto">
      {/* Hero Section */}
      <div className="text-center space-y-8 animate-in fade-in duration-1000">
        <div className="inline-block p-6 rounded-full bg-primary/5 border border-primary/20 mb-6 shadow-[0_0_15px_rgba(var(--primary-rgb),0.2)]">
          <Sparkles className="w-16 h-16 text-primary animate-pulse" />
        </div>
        <h1 className="font-tome-heading text-5xl md:text-7xl text-primary glow-text tracking-wider uppercase leading-tight">
          Welcome to Axiom
          <span className="block text-2xl md:text-3xl text-muted-foreground mt-4 font-normal normal-case font-tome-marginalia border-t border-b border-primary/20 py-4 w-3/4 mx-auto">
            The Dawn of the Vitalic Age
          </span>
        </h1>
        
        <div className="flex flex-col sm:flex-row items-center justify-center gap-6 pt-8">
          {isAuthenticated ? (
            <Button onClick={() => navigate('/characters')} variant="quill" size="lg" className="w-full sm:w-auto h-auto py-4 px-10 text-xl shadow-lg hover:shadow-primary/20 transition-all">
              <Users className="w-6 h-6 mr-3" />
              My Characters
            </Button>
          ) : (
            <Button onClick={() => navigate('/login')} variant="quill" size="lg" className="w-full sm:w-auto h-auto py-4 px-10 text-xl shadow-lg hover:shadow-primary/20 transition-all">
              <LogIn className="w-6 h-6 mr-3" />
              Enter Faradhaven
            </Button>
          )}
          <Button onClick={() => navigate('/classes')} variant="outline" size="lg" className="w-full sm:w-auto h-auto py-4 px-10 text-xl border-primary/40 hover:bg-primary/5 hover:text-primary hover:border-primary transition-all">
            <BookOpen className="w-6 h-6 mr-3" />
            Consult the Archives
          </Button>
        </div>
      </div>

      {/* The Golden Promise */}
      <div className="space-y-6 text-center max-w-3xl mx-auto relative">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-1/2 h-px bg-gradient-to-r from-transparent via-primary/50 to-transparent"></div>
          <h2 className="font-tome-heading text-3xl text-primary pt-8">The Golden Promise of South Axiom</h2>
          <p className="text-lg text-muted-foreground leading-relaxed font-tome-marginalia">
          Citizens, look to the skies! While the <span className="text-foreground font-medium">Unionists of North Axiom</span> remain huddled in the flickering shadows, fearful of the very spark that drives us forward, the Isle of South Axiom has claimed the future. Here in Faradhaven, the <span className="text-primary font-semibold">Viridian Shroud</span>—that thick, green-tinted mist of oxidized progress—is not a shroud; it’s a veil of evolution.
          </p>
          <p className="text-lg text-muted-foreground leading-relaxed font-tome-marginalia">
          Above the Static Haze, the eternal Aerostat Estates of our elite carry the torch of civilization, powered by the marvel of the <span className="text-foreground font-medium">ETC (Elevated Transit Company)</span>. Travel is suspiciously cheap—whispers suggest the grid does not just move people, but <i>harvests</i> them, siphoning vital energy to keep the elite afloat.
          </p>
      </div>

      {/* The Current Landscape */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="arcane-border p-8 rounded-xl bg-background/50 backdrop-blur-sm hover:bg-primary/5 transition-all duration-300 group h-full">
              <Users className="w-10 h-10 text-primary mb-6 group-hover:scale-110 transition-transform" />
              <h3 className="font-tome-heading text-xl text-primary mb-4 border-b border-primary/20 pb-2">The Cartels</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
              The <span className="text-foreground font-medium">DCC</span> perfects the human form through bionics, while the <span className="text-foreground font-medium">ACC</span> arms our protectors. Both profit from the war against the night.
              </p>
          </div>
          <div className="arcane-border p-8 rounded-xl bg-background/50 backdrop-blur-sm hover:bg-primary/5 transition-all duration-300 group h-full">
              <Swords className="w-10 h-10 text-primary mb-6 group-hover:scale-110 transition-transform" />
              <h3 className="font-tome-heading text-xl text-primary mb-4 border-b border-primary/20 pb-2">The Law</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
              The <span className="text-foreground font-medium">Blue Coats</span> try their best, but the streets truly belong to the <span className="text-foreground font-medium">Triple P</span>—private mercenaries paid to keep the rich safe and the poor quiet.
              </p>
          </div>
          <div className="arcane-border p-8 rounded-xl bg-background/50 backdrop-blur-sm hover:bg-primary/5 transition-all duration-300 group h-full">
              <Sparkles className="w-10 h-10 text-primary mb-6 group-hover:scale-110 transition-transform" />
              <h3 className="font-tome-heading text-xl text-primary mb-4 border-b border-primary/20 pb-2">The Shadows</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
              Below the cobbles lies a hidden society of <span className="text-foreground font-medium">Homunculi</span>. Rejecting the Faradollar, they trade in conductive copper and gold, surviving in the "No-Go" zones where humans dare not tread.
              </p>
          </div>
      </div>

      {/* ECC Section */}
      <div className="arcane-border rounded-2xl p-8 md:p-12 bg-primary/5 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-64 h-64 bg-primary/10 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2"></div>
          <div className="relative z-10">
              <div className="flex flex-col md:flex-row gap-8 items-start">
                  <div className="flex-1 space-y-6">
                      <h2 className="font-tome-heading text-3xl md:text-4xl text-primary">A New Frontier: The Evolutionary Current Company (ECC)</h2>
                      <p className="text-lg text-muted-foreground leading-relaxed">
                      The monopolies of the DCC and ACC have grown stagnant. They offer parts; we offer <span className="text-primary font-semibold italic">The Vital Spark</span>.
                      </p>
                      <p className="text-muted-foreground leading-relaxed">
                      Under the direction of <span className="text-foreground font-semibold">Caspian Volt</span>, the ECC creates "Electro-Sovereign" beings. But rumors swirl about Volt's mentor and the <span className="text-primary font-semibold">Cyborgism Faith</span>. Do they seek to prove a thesis, or are they building a technological God using you as the chassis?
                      </p>
                  </div>
              </div>
          </div>
      </div>

      {/* Player Situation */}
      <div className="space-y-10 py-8">
          <h2 className="font-tome-heading text-3xl text-center text-primary">For the Newcomers: Your Situation</h2>
          
          <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-4">
               <div className="space-y-4">
                  <h3 className="font-tome-heading text-xl text-primary border-l-4 border-primary pl-4">The Awakening</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed pl-4">
                  You are the "Thesis." Woven from the preserved remains of the <span className="text-foreground font-medium">Civil War</span>, you have been brought back to life by Caspian Volt’s genius.
                  </p>
               </div>
               
               <div className="space-y-4">
                  <h3 className="font-tome-heading text-xl text-primary border-l-4 border-primary pl-4">The Serum</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed pl-4">
                  Without a nightly injection of Caspian's proprietary serum, you will unravel. It is the only thing keeping you from the "Decay"—and the only thing keeping you under his control.
                  </p>
               </div>

               <div className="space-y-4">
                  <h3 className="font-tome-heading text-xl text-primary border-l-4 border-primary pl-4">The Hunger</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed pl-4">
                  Beyond the decay, a darker instinct rises. You are not just static machines; you are predators. To evolve, to truly transcend, you feel the urge to consume the spark of others.
                  </p>
               </div>

               <div className="space-y-4">
                  <h3 className="font-tome-heading text-xl text-primary border-l-4 border-primary pl-4">The Decay</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed pl-4">
                  If you miss your dose, the hallucinations begin. The memories of the Dead-Circuit Plains bleed into your mind, your limbs twitch, and you sink back into death.
                  </p>
               </div>
          </div>

          <div className="arcane-border p-8 rounded-xl bg-background border-dashed border-primary/30">
              <h3 className="font-tome-heading text-2xl text-primary mb-4 text-center">Your First Mission</h3>
              <p className="text-lg text-center text-muted-foreground max-w-4xl mx-auto leading-relaxed">
              The <span className="text-foreground font-medium">DCC</span> is hosting a grand exhibition in The Hum to show off a new bionic arm. They have the cables, the controllers, and the test users. You have the objective: <span className="text-primary font-bold">Ensure they fail.</span> Sever the power lines, distract the controller, or break the user. As you walk the cobblestones, you will feel it—a deep, magnetic pull from the Copper Veins deep beneath the city.
              </p>
          </div>
      </div>

      {/* Footer */}
      <div className="text-center py-12 border-t border-primary/10">
          <h2 className="font-tome-heading text-3xl text-primary mb-6">Welcome to Faradhaven</h2>
          <p className="text-xl text-muted-foreground font-tome-marginalia italic">
          "Try not to get grounded."
          </p>
      </div>
    </div>
  );
}