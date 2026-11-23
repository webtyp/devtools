# Bugs de UI y mejoras — `webtyp/devtools`

> Auditoría en vivo sobre `#api_explorer` y `#db_schema` (2026-11-19).

## Bugs encontrados

| # | Sev | Pantalla | Qué se ve en pantalla | Causa raíz verificada |
|---|---|---|---|---|
| DT-1 | 🟡 | `#api_explorer` | **Botón «Enviar» sin estilo visual de botón.** Al expandir una ruta (`GET /__webtyp/ca`), el texto `Enviar` flota como texto plano sin fondo, borde ni affordance clicable sobre una caja blanca vacía. | `webtyp/layout/apiexplorer` no aplica receta de botón (`style.Button(...)`) al disparador de prueba de ruta. |
| DT-2 | 🟡 | `#api_explorer` | **Etiquetas de política en inglés** (`public`, `guarded`, `specialty:c`, `service:r`) junto a UI en español (`Buscar rutas`, `Enviar`). | Strings literales sin `lang.Text` / `lang.json` en `webtyp/layout/apiexplorer`. |
| DT-3 | 🟡 | `#db_schema` | **18 avisos de integridad en el diagrama de BD** (`0 errores · 18 avisos`): `user_role` aparece sin líneas FK hacia `user`/`role`; `access_audit`, `access_record_link`, `business_hours`, `client`, `closure`, `conflict_recompute` aparecen como tablas aisladas. | `rbac` (`user_role`) y tablas de auditoría/calendario no declaran `Ref` o `schemaview` marca como aviso tablas legítimamente independientes de catálogo raíz (`business_hours`, `closure`, `client`). |
| DT-4 | ⚪ | Menú `devtools` | **`ViewAs` configurado en `devtools/menu.go:18` pero no visible ni accesible en la UI.** | `devtools.Menu()` declara `ViewAs: &appmenu.ViewAsConfig{}`, pero en la barra lateral bajo `Herramientas de desarrollo` solo aparecen `Explorador de API` y `Base de datos`. |

## Sugerencias de mejora justificadas

1. **Estilizar el botón «Enviar» y traducir `public`/`guarded` en `apiexplorer`:** aplicar `style.Button(style.Primary)` al botón de envío y pasar los estados de acceso a `lang.Text` (`Pública`, `Protegida`). *Justificación:* coherencia visual con `ui-consistency` y claridad inmediata de qué botón ejecuta la petición.
2. **Declarar FKs faltantes en `access` (`user_role`) y distinguir tablas raíz en `schemaview`:** conectar `user_role` con `user` y `role` (como ya reportó `AUDITORIA_MVP.md` H11) y permitir marcar tablas de configuración global/calendario sin falso aviso de «sin relación». *Justificación:* evita ruido permanente («18 avisos») que oculta avisos reales de modelado cuando se agregan módulos nuevos.
