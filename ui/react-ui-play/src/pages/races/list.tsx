import {
    List,
    Table,
    ShowButton,
    useTable,
} from "@pankod/refine-antd";

import { IRace } from "interfaces";

export const RaceList: React.FC = () => {
    const { tableProps } = useTable<IRace>();

    return (
        <List>
            <Table {...tableProps} rowKey="id">
                <Table.Column dataIndex="name" title="Name" />
                <Table.Column dataIndex="description" title="Description" />
                <Table.Column dataIndex={["base_attributes", "strength"]} title="Strength"/>
                <Table.Column dataIndex={["base_attributes", "agility"]} title="Agility" />
                <Table.Column dataIndex={["base_attributes", "intelligence"]} title="Intelligence"/>
                <Table.Column dataIndex={["base_attributes", "willpower"]} title="Willpower"/>
                <Table.Column<IRace>
                    title="Actions"
                    dataIndex="actions"
                    render={(_text, record): React.ReactNode => {
                        return (
                            <ShowButton
                                size="small"
                                recordItemId={record.id}
                                hideText
                            />
                        );
                    }}
                />
            </Table>
        </List>
    );
};