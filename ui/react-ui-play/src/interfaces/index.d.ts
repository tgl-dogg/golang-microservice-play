import internal from "stream";

export interface IRace {
    id: string;
    name: string;
    description: string;
    base_attributes: BaseAttributes;
}

interface BaseAttributes {
    strength: int;
    agility: int;
    intelligence: int;
    willpower: internal;
}