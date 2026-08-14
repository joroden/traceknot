import { useEffect, useState } from "react";
import { fetchNodeDetail, type NodeDetail } from "../api";

export interface NodeDetailState {
  detail: NodeDetail | null;
  launchDetail: NodeDetail | null;
  loading: boolean;
  error: string | null;
}

export function useNodeDetail(
  nodeId: string | null,
  launchNodeId: string | null,
): NodeDetailState {
  const [detail, setDetail] = useState<NodeDetail | null>(null);
  const [launchDetail, setLaunchDetail] = useState<NodeDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!nodeId) {
      setDetail(null);
      setLaunchDetail(null);
      setError(null);
      setLoading(false);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(null);
    Promise.all([
      fetchNodeDetail(nodeId),
      launchNodeId ? fetchNodeDetail(launchNodeId) : Promise.resolve(null),
    ])
      .then(([node, launch]) => {
        if (!cancelled) {
          setDetail(node);
          setLaunchDetail(launch);
          setLoading(false);
        }
      })
      .catch((reason: unknown) => {
        if (!cancelled) {
          setError(reason instanceof Error ? reason.message : String(reason));
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [nodeId, launchNodeId]);

  return { detail, launchDetail, loading, error };
}
