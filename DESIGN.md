# FlagBridge Design System

> Architectural Precision — "The Precision Engineer"

Design system para o FlagBridge, plataforma de gerenciamento de feature flags. A UI funciona como um instrumento de alta performance — como uma IDE premium ou uma Leica. O princípio central é **Subtractive Design**: remover tudo que não é essencial para deixar os dados e a lógica das flags respirarem.

## Fundamentos

- **Color Mode:** Dark
- **Primary Color:** Electric Blue `#3B82F6`
- **Secondary Color:** Slate `#64748B`
- **Neutral Base:** `#0F172A`
- **Font Stack:** Inter (body, headline, label)
- **Corner Radius:** 0.375rem (md) padrão
- **Spacing Scale:** 1x (base)
- **Target:** Desktop-first (1280px), com suporte mobile (390px)

---

## Cores e Arquitetura Tonal

### Superficies (do mais profundo ao mais elevado)

| Token                       | Hex       | Uso                                      |
| --------------------------- | --------- | ---------------------------------------- |
| `surface_container_lowest`  | `#060e20` | Base absoluta, nunca usar `#000`         |
| `surface` / `surface_dim`   | `#0b1326` | Background principal                     |
| `surface_container_low`     | `#131b2e` | Sidebar, seções secundárias              |
| `surface_container`         | `#171f33` | Cards, containers padrão                 |
| `surface_container_high`    | `#222a3d` | Área de foco ativa, workspaces           |
| `surface_container_highest` | `#2d3449` | Modais, dropdowns, overlays              |
| `surface_bright`            | `#31394d` | Superfície de destaque                   |

### Regra "No-Line"

Bordas de 1px são proibidas para separação de seções. Usar transições de cor entre superfícies para definir limites visuais.

- Sidebar `surface_container_low` contra main `surface`
- A borda é definida pela transição de cor, não por um stroke

### Regra "Glass & Gradient"

- CTAs principais: gradiente de `primary` (`#adc6ff`) para `primary_container` (`#4d8eff`) a 135deg
- Navegação flutuante / command palette: `surface_container_highest` com `backdrop-blur(20px)` (efeito "frosted sapphire")

### Paleta Completa

**Primary**
| Token                    | Hex       |
| ------------------------ | --------- |
| `primary`                | `#adc6ff` |
| `primary_container`      | `#4d8eff` |
| `primary_fixed`          | `#d8e2ff` |
| `primary_fixed_dim`      | `#adc6ff` |
| `on_primary`             | `#002e6a` |
| `on_primary_container`   | `#00285d` |
| `on_primary_fixed`       | `#001a42` |
| `inverse_primary`        | `#005ac2` |

**Secondary**
| Token                    | Hex       |
| ------------------------ | --------- |
| `secondary`              | `#b7c8e1` |
| `secondary_container`    | `#3a4a5f` |
| `secondary_fixed`        | `#d3e4fe` |
| `secondary_fixed_dim`    | `#b7c8e1` |
| `on_secondary`           | `#213145` |
| `on_secondary_container` | `#a9bad3` |

**Tertiary (Warnings)**
| Token                     | Hex       |
| ------------------------- | --------- |
| `tertiary`                | `#ffb786` |
| `tertiary_container`      | `#df7412` |
| `tertiary_fixed`          | `#ffdcc6` |
| `tertiary_fixed_dim`      | `#ffb786` |
| `on_tertiary`             | `#502400` |
| `on_tertiary_container`   | `#461f00` |

**Error**
| Token              | Hex       |
| ------------------ | --------- |
| `error`            | `#ffb4ab` |
| `error_container`  | `#93000a` |
| `on_error`         | `#690005` |
| `on_error_container` | `#ffdad6` |

**Text & Outline**
| Token                | Hex       |
| -------------------- | --------- |
| `on_surface`         | `#dae2fd` |
| `on_surface_variant`  | `#c2c6d6` |
| `on_background`      | `#dae2fd` |
| `outline`            | `#8c909f` |
| `outline_variant`    | `#424754` |
| `inverse_surface`    | `#dae2fd` |
| `inverse_on_surface` | `#283044` |
| `surface_tint`       | `#adc6ff` |

---

## Tipografia

Font: **Inter** para todos os níveis.

| Token        | Tamanho   | Uso                                            |
| ------------ | --------- | ---------------------------------------------- |
| `display-md` | 2.75rem   | Hero states, letter-spacing: -0.02em           |
| `title-lg`   | 1.375rem  | Nomes de flags — transmite autoridade          |
| `body-md`    | 0.875rem  | Texto principal, line-height: 1.5              |
| `label-md`   | —         | Labels de inputs, cor `on_surface_variant`     |
| `label-sm`   | 0.6875rem | Metadata: tags de ambiente, timestamps         |

---

## Elevacao e Profundidade

Profundidade via luminosidade, nao via bordas.

- **Layering:** empilhar containers para criar lift natural (`surface_container_lowest` sobre `surface_container_low` = efeito inset)
- **Ambient Shadows (floating):** `0 20px 40px rgba(6, 14, 32, 0.4)` — cor da sombra sempre tintada com `surface_container_lowest`
- **Ghost Border (fallback):** se uma borda for absolutamente necessária (ex: tabelas densas), usar `outline_variant` (`#424754`) a 15% de opacidade
- **Glassmorphism:** navbars com `surface` a 80% opacidade + `blur(12px)`

---

## Componentes

### Buttons

| Variante  | Background                                        | Texto           | Radius |
| --------- | ------------------------------------------------- | --------------- | ------ |
| Primary   | Gradiente `primary` -> `primary_container` (135deg) | `on_primary`    | md     |
| Tertiary  | Transparente (hover: `surface_container_high`)    | `primary`       | md     |

### Inputs

- Background: `surface_container_highest`
- Focus: 1px `primary` ghost border + glow `primary` (4px spread)
- Labels: `label-md` em `on_surface_variant`, sempre acima do input

### Cards & Lists

- Cards: `surface_container` com radius `lg` (0.5rem)
- Separação entre itens: padding vertical de 0.6–0.9rem, sem dividers
- Se divider necessário: 1px `surface_variant` cobrindo 90% da largura (estilo editorial indentado)

### Feature Flag Toggle (componente custom)

| Estado   | Track                          | Thumb                  |
| -------- | ------------------------------ | ---------------------- |
| Inativo  | `secondary_container` (#3a4a5f) | `secondary_fixed` (#d3e4fe) |
| Ativo    | `primary_container` (#4d8eff)  | `primary_fixed` (#d8e2ff)   |

Animação: spring 0.2s (snappy, como switch físico)

---

## Espacamento

- Margens top-level: `4.5rem` (20) ou `5.5rem` (24) para sensação "premium"
- Se um layout parecer apertado, aumentar padding em 2 steps na escala

---

## Regras

**Fazer:**
- Usar transições tonais em vez de bordas
- Usar `on_tertiary_container` para warnings (burnt orange sofisticado)
- Abraçar whitespace generoso

**Nao fazer:**
- Bordas opacas de alto contraste
- Preto puro (`#000000`) — o mais escuro é `surface_container_lowest` (`#060e20`)
- Drop shadows padrão em cards
- Texto longo centralizado — manter grid left-aligned, developer-first
