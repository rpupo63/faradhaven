import { ApiBeast } from '@/types/game/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

interface BeastCardProps {
  beast: ApiBeast;
}

export function BeastCard({ beast }: BeastCardProps) {
  return (
    <Card className="flex flex-col h-full min-w-0 overflow-hidden">
      <CardHeader className="min-w-0">
        <CardTitle className="flex flex-col gap-2 min-[420px]:flex-row min-[420px]:justify-between min-[420px]:items-center text-primary glow-text">
          <span className="break-words min-w-0">{beast.name}</span>
          <Badge variant="secondary" className="shrink-0 w-fit">CR: {beast.challenge_rating}</Badge>
        </CardTitle>
        <div className="flex flex-wrap gap-2 mt-2">
          <Badge variant="outline">{beast.size}</Badge>
          <Badge variant="outline">{beast.type}</Badge>
          <Badge variant="outline">{beast.alignment}</Badge>
        </div>
      </CardHeader>
      <CardContent className="flex-grow space-y-2 text-sm min-w-0">
        {beast.image_url && (
          <img src={beast.image_url} alt={beast.name} className="w-full h-32 object-cover rounded-md mb-2" />
        )}
        <p><strong>AC:</strong> {beast.armor_class}</p>
        <p><strong>HP:</strong> {beast.hit_points} ({beast.hit_dice})</p>
        <p><strong>Speed:</strong> {beast.speed}</p>
        <div className="grid grid-cols-2 gap-1 min-[380px]:grid-cols-3">
          <p><strong>STR:</strong> {beast.strength}</p>
          <p><strong>DEX:</strong> {beast.dexterity}</p>
          <p><strong>CON:</strong> {beast.constitution}</p>
          <p><strong>INT:</strong> {beast.intelligence}</p>
          <p><strong>WIS:</strong> {beast.wisdom}</p>
          <p><strong>CHA:</strong> {beast.charisma}</p>
        </div>
        {beast.description && <p className="mt-2 text-muted-foreground line-clamp-3">{beast.description}</p>}
        {beast.abilities && beast.abilities.length > 0 && (
          <div className="break-words">
            <strong>Abilities:</strong> <span className="text-muted-foreground">{beast.abilities.join(', ')}</span>
          </div>
        )}
        {beast.actions && beast.actions.length > 0 && (
          <div className="break-words">
            <strong>Actions:</strong> <span className="text-muted-foreground">{beast.actions.join(', ')}</span>
          </div>
        )}
        {beast.attacks && beast.attacks.length > 0 && (
          <div>
            <strong>Attacks:</strong>
            <ul className="list-disc list-inside text-muted-foreground">
              {beast.attacks.map((attack) => (
                <li key={attack.id}>
                  {attack.name} (+{attack.attack_bonus}) {attack.damage_dice} {attack.damage_type}
                </li>
              ))}
            </ul>
          </div>
        )}
        {beast.skills && beast.skills.length > 0 && (
          <div>
            <strong>Skills:</strong>
            <ul className="list-disc list-inside text-muted-foreground">
              {beast.skills.map((skill) => (
                <li key={skill.id}>
                  {skill.name} ({skill.modifier})
                </li>
              ))}
            </ul>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
