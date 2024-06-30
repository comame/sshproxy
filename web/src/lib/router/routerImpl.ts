import { ComponentType } from "svelte";
import { writable } from "svelte/store";

export type routesDef = {
  [path: string]: ComponentType;
};

type currentRouteStore = [ComponentType];
export const currentRoute = writable<currentRouteStore | null>(null);

type setupRouterOptions = {
  routes: routesDef;
  hashRouter: boolean;
  fallback: ComponentType | null;
};

export function setupRouter(window: Window, options: setupRouterOptions) {
  if (!options.hashRouter) {
    throw "hashRouter 以外はサポートできてない";
  }

  // キーにパスを判定する正規表現を含むオブジェクト
  const regexRoutes: routesDef = {};
  for (const path of Object.keys(options.routes)) {
    const r = convertPathToRegex(path);
    regexRoutes[r] = options.routes[path];
  }

  const routeChanged = (newPath: string) => {
    const r = resolve(newPath, regexRoutes);
    if (!r) {
      if (options.fallback) {
        currentRoute.set([options.fallback]);
      } else {
        currentRoute.set(null);
      }
      return;
    }

    currentRoute.set([r]);
  };

  const hashChange = () => {
    routeChanged(hashToPath(window.location.hash));
  };
  window.addEventListener("hashchange", hashChange);

  if (options.hashRouter) {
    // 初回は hashchange が呼ばれないので手動で呼ぶ
    hashChange();
  }

  return {
    destroy() {
      window.removeEventListener("hashchange", hashChange);
    },
  };
}

function hashToPath(hash: string): string {
  // #/abc => /abc
  // '' => /
  return "/" + hash.slice(2);
}

function convertPathToRegex(path: string): string {
  const paths = path.split("/");
  const rs: string[] = [];
  for (const p of paths) {
    if (p.startsWith(":")) {
      throw "ダイナミックルートはまだ非対応";

      // FIXME: ダイナミックルート部分にこれら以外の文字列が来たとき正常にルーティングできない
      rs.push("[a-zA-Z0-9-_]");
      continue;
    }
    // FIXME: 正規表現をエスケープしていないので、怪しいルート定義が渡されると壊れる
    rs.push(p);
  }
  return "^" + rs.join("\\/") + "$";
}

// ユーザーが渡したパスとキーにパスを判定する正規表現を含むオブジェクトを受け取って、レンダリングするべきコンポーネントを返す
// FIXME: O(N * <正規表現の実行>) かかるのでパフォーマンスは良くない
function resolve(path: string, def: routesDef): ComponentType | null {
  for (const key of Object.keys(def)) {
    const match = new RegExp(key).exec(path);
    if (match === null) {
      continue;
    }

    return def[key];
  }

  return null;
}
